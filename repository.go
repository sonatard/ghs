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
	opts       *github.SearchOptions
	query      string
	repos      []github.Repository
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
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		Debug("Github token get from environment value\n")
		return token
	}

	// github.token in gitconfig
	token, err := gitconfig.GetString("github.token")
	if err == nil {
		Debug("Github token get from gitconfig value\n")
		return token
	}

	Debug("Github token not found\n")
	return ""
}

func NewRepo(sort string, order string, max int, enterprise string, token string, query string) (*repo, error) {

	// Github Token authentication
	githubToken := getToken(token)

	var tc *http.Client
	if githubToken != "" {
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

	// Github API Search options
	searchOpts := &github.SearchOptions{
		Sort:        sort,
		Order:       order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	return &repo{client: cli, opts: searchOpts, query: query, printMax: max, printCount: 0}, nil
}

func (r *repo) search() (repos []github.Repository) {
	Debug("%d go func search start\n", r.opts.ListOptions.Page)
	ret, _, err := r.client.Search.Repositories(r.query, r.opts)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", r.query)
		fmt.Println(err)
	}
	Debug("%d go func search end\n", r.opts.ListOptions.Page)

	return ret.Repositories
}

func (r *repo) SearchRepository() (<-chan []github.Repository, <-chan bool) {
	var wg sync.WaitGroup
	reposBuff := make(chan []github.Repository, 10)
	fin := make(chan bool)

	ret, resp, err := r.client.Search.Repositories(r.query, r.opts)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", r.query)
		fmt.Println(err)
		os.Exit(1)
	}
	reposBuff <- ret.Repositories
	last := resp.LastPage
	Debug("LastPage = %d\n", last)

	go func() {
		for i := 1; i < last; i++ {
			Debug("main thread %d\n", i)
			wg.Add(1)
			r.opts.ListOptions.Page = i
			go func() {
				reposBuff <- r.search()
				wg.Done()
			}()
		}
		Debug("main thread wait...\n")
		wg.Wait()
		Debug("main thread wakeup!!\n")
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
