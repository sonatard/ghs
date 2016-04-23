package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/motemen/go-gitconfig"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type SearchInfo struct {
	sort    string
	order   string
	query   string
	max     int
	perPage int
}

type repo struct {
	client     *github.Client
	info       *SearchInfo
	printCount int
}

func getToken(optsToken string) string {
	// -t or --token option
	if optsToken != "" {
		Debug("Github token get from option value\n")
		return optsToken
	}

	// GITHUB_TOKEN environment
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		Debug("Github token get from environment value\n")
		return token
	}

	// github.token in gitconfig
	if token, err := gitconfig.GetString("github.token"); err == nil {
		Debug("Github token get from gitconfig value\n")
		return token
	}

	Debug("Github token not found\n")
	return ""
}

func NewRepo(info *SearchInfo, baseURL *url.URL, token string) (*repo, error) {
	var tc *http.Client

	if githubToken := getToken(token); githubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	cli := github.NewClient(tc)

	if baseURL != nil {
		cli.BaseURL = baseURL
	}

	return &repo{
		client:     cli,
		info:       info,
		printCount: 0}, nil
}

func (r *repo) search(page int) (repos []github.Repository) {
	Debug("Page%d go func search start\n", page)

	opts := &github.SearchOptions{
		Sort:        r.info.sort,
		Order:       r.info.order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: r.info.perPage, Page: page},
	}

	Debug("Page%d query : %s\n", page, r.info.query)
	ret, _, err := r.client.Search.Repositories(r.info.query, opts)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", r.info.query)
		fmt.Println(err)
	}
	Debug("Page%d go func search result length %d\n", page, len(ret.Repositories))
	Debug("Page%d go func search end\n", page)

	return ret.Repositories
}

func (r *repo) SearchRepository() (<-chan []github.Repository, <-chan bool) {
	var wg sync.WaitGroup
	reposBuff := make(chan []github.Repository, 10)
	fin := make(chan bool)
	page := 0

	opts := &github.SearchOptions{
		Sort:        r.info.sort,
		Order:       r.info.order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: r.info.perPage, Page: page},
	}
	Debug("Page%d query : %s\n", page, r.info.query)
	ret, resp, err := r.client.Search.Repositories(r.info.query, opts)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", r.info.query)
		fmt.Println(err)
		os.Exit(1)
	}
	Debug("Total = %d\n", ret.Total)
	reposBuff <- ret.Repositories
	Debug("main thread repos length %d\n", len(ret.Repositories))
	last := ((r.info.max - 1) / r.info.perPage)
	if resp.LastPage < last {
		last = resp.LastPage
	}

	Debug("resp.LastPage = %d\n", resp.LastPage)
	Debug("LastPage = %d\n", last)
	page++

	go func() {
		for ; page < last+1; page++ {
			Debug("sub thread %d\n", page)
			wg.Add(1)
			go func(p int) {
				reposBuff <- r.search(p)
				wg.Done()
			}(page)
		}
		Debug("sub thread wait...\n")
		wg.Wait()
		Debug("sub thread wakeup!!\n")
		fin <- true
	}()

	Debug("main thread return\n")

	return reposBuff, fin
}

func (r *repo) PrintRepository(repos []github.Repository) (end bool) {
	Debug("repos length %d\n", len(repos))
	repoNameMaxLen := 0
	for _, repo := range repos {
		repoNamelen := len(*repo.FullName)
		if repoNamelen > repoNameMaxLen {
			repoNameMaxLen = repoNamelen
		}
	}
	for _, repo := range repos {
		if repo.FullName != nil {
			printf("%v", *repo.FullName)
		}

		printf("    ")

		paddingLen := repoNameMaxLen - len(*repo.FullName)

		for i := 0; i < paddingLen; i++ {
			printf(" ")
		}

		if repo.Description != nil {
			printf("%v", *repo.Description)
		}

		printf("\n")

		r.printCount++
		Debug("printCount %d, max %d\n", r.printCount, r.info.max)
		if r.printCount >= r.info.max {
			return true
		}

	}
	return false
}

func printf(format string, args ...interface{}) {
	if os.Getenv("GHS_PRINT") != "no" {
		fmt.Printf(format, args...)
	}
}
