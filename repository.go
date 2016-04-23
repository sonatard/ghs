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
	client      *github.Client
	info        *SearchInfo
	last_page   int
	max_count   int
	print_count int
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
		client:      cli,
		info:        info,
		print_count: 0}, nil
}

func s(client *github.Client, page int, si *SearchInfo) (*github.RepositoriesSearchResult, *github.Response, error) {
	opts := &github.SearchOptions{
		Sort:        si.sort,
		Order:       si.order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: si.perPage, Page: page},
	}
	return client.Search.Repositories(si.query, opts)
}

func (r *repo) first_search() (repos []github.Repository, last_page int, max_count int) {
	ret, resp, err := s(r.client, 1, r.info)

	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", r.info.query)
		fmt.Println(err)
		os.Exit(1)
	}
	Debug("main thread repos length %d\n", len(ret.Repositories))

	max := r.info.max
	Debug("Total = %d\n", *ret.Total)
	Debug("r.info.max = %d\n", r.info.max)
	if *ret.Total < max {
		max = *ret.Total
	}

	last := ((max - 1) / r.info.perPage) + 1
	Debug("resp.LastPage = %d\n", resp.LastPage)
	Debug("last = %d\n", last)
	if resp.LastPage < last {
		last = resp.LastPage
	}

	return ret.Repositories, last, max
}

func (r *repo) search(page int) (repos []github.Repository) {
	Debug("Page%d go func search start\n", page)

	Debug("Page%d query : %s\n", page, r.info.query)
	ret, _, err := s(r.client, page, r.info)
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
	reposBuff := make(chan []github.Repository, 1000)
	fin := make(chan bool)

	// 1st search
	repos, last, max := r.first_search()

	// notify main thread of first search result
	r.last_page = last
	r.max_count = max
	reposBuff <- repos
	Debug("Search settings r.last_page = %d\n", r.last_page)
	Debug("Search settings r.max_count = %d\n", r.max_count)

	// 2nd - 10th search
	go func() {
		for page := 2; page < last+1; page++ {
			Debug("sub thread %d\n", page)
			wg.Add(1)
			go func(p int) {
				// notify main thread of 2nd - 10th search result
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

		r.print_count++
		Debug("print_count %d, max_count %d\n", r.print_count, r.max_count)
		if r.print_count >= r.max_count {
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
