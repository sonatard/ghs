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

type repo struct {
	client     *github.Client
	sort       string
	order      string
	query      string
	repos      []github.Repository
	perPage    int
	printMax   int
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

func NewRepo(sort string, order string, max int, enterprise string, token string, query string) (*repo, error) {
	var tc *http.Client

	if githubToken := getToken(token); githubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	cli := github.NewClient(tc)

	// Github API
	if enterprise != "" {
		baseURL, err := url.Parse(enterprise)
		if err != nil {
			return nil, err
		}
		cli.BaseURL = baseURL
	}

	return &repo{
		client:     cli,
		sort:       sort,
		order:      order,
		query:      query,
		perPage:    100,
		printMax:   max,
		printCount: 0}, nil
}

func (r *repo) search(page int) (repos []github.Repository) {
	Debug("Page%d go func search start\n", page)

	opts := &github.SearchOptions{
		Sort:        r.sort,
		Order:       r.order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: r.perPage, Page: page},
	}

	Debug("Page%d query : %s\n", page, r.query)
	ret, _, err := r.client.Search.Repositories(r.query, opts)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", r.query)
		fmt.Println(err)
	}
	Debug("Page%d go func search end\n", page)

	return ret.Repositories
}

func (r *repo) SearchRepository() (<-chan []github.Repository, <-chan bool) {
	var wg sync.WaitGroup
	reposBuff := make(chan []github.Repository, 10)
	fin := make(chan bool)
	page := 0

	opts := &github.SearchOptions{
		Sort:        r.sort,
		Order:       r.order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: r.perPage, Page: page},
	}
	Debug("Page%d query : %s\n", page, r.query)
	ret, resp, err := r.client.Search.Repositories(r.query, opts)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", r.query)
		fmt.Println(err)
		os.Exit(1)
	}
	Debug("Total = %d\n", ret.Total)
	reposBuff <- ret.Repositories
	Debug("main thread repos length %d\n", len(ret.Repositories))
	last := ((r.printMax - 1) / r.perPage)
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

func (r *repo) PrintRepository() (end bool) {
	Debug("r.repos length %d\n", len(r.repos))
	repoNameMaxLen := 0
	for _, repo := range r.repos {
		repoNamelen := len(*repo.FullName)
		if repoNamelen > repoNameMaxLen {
			repoNameMaxLen = repoNamelen
		}
	}
	for _, repo := range r.repos {
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
		Debug("printCount %d\n", r.printCount)
		if r.printCount >= r.printMax {
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
