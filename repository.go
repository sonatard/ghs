package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"net/url"
	"sync"
)

type repo struct {
	client *github.Client
	opts   *github.SearchOptions
	query  string
	repos  []github.Repository
}

func NewRepo(sort string, order string, max int, enterprise string, query string) (*repo, error) {
	var r repo

	client := github.NewClient(nil)

	if enterprise != "" {
		baseURL, err := url.Parse(enterprise)
		if err != nil {
			return nil, err
		}
		client.BaseURL = baseURL
	}

	searchOpts := &github.SearchOptions{
		Sort:        sort,
		Order:       order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	r.client = client
	r.opts = searchOpts
	r.query = query

	return &r, nil
}

func (r repo) search() (repos []github.Repository) {
	fmt.Printf("%d go func search start\n", r.opts.ListOptions.Page)
	ret, _, err := r.client.Search.Repositories(r.query, r.opts)
	if err != nil {
		fmt.Printf("Search Error %s\n", r.query)
	}
	fmt.Printf("%d go func search end\n", r.opts.ListOptions.Page)

	return ret.Repositories
}

func (r repo) SearchRepository() (<-chan []github.Repository, <-chan bool) {
	var wg sync.WaitGroup
	reposBuff := make(chan []github.Repository, 10)
	fin := make(chan bool)

	ret, resp, err := r.client.Search.Repositories(r.query, r.opts)
	if err != nil {
		fmt.Printf("Search Error %s\n", r.query)
	}
	reposBuff <- ret.Repositories
	last := resp.LastPage
	fmt.Printf("LastPage = %d\n", last)

	go func() {
		for i := 0; i < last; i++ {
			fmt.Printf("main thread %d\n", i)
			wg.Add(1)
			r.opts.ListOptions.Page = i
			go func() {
				reposBuff <- r.search()
				wg.Done()
			}()
		}
		fmt.Printf("main thread wait...\n")
		wg.Wait()
		fmt.Printf("main thread wakeup!!\n")
		fin <- false
	}()

	fmt.Printf("main thread return\n")

	return reposBuff, fin
}

func (r repo) PrintRepository() {
	fmt.Printf("%d\n", len(r.repos))
	repoNameMaxLen := 0
	for _, repo := range r.repos {
		repoNamelen := len(*repo.FullName)
		if repoNamelen > repoNameMaxLen {
			repoNameMaxLen = repoNamelen
		}
	}
	for _, repo := range r.repos {
		if repo.FullName != nil {
			fmt.Printf("%v", *repo.FullName)
		}

		fmt.Printf("    ")

		paddingLen := repoNameMaxLen - len(*repo.FullName)

		for i := 0; i < paddingLen; i++ {
			fmt.Printf(" ")
		}

		if repo.Description != nil {
			fmt.Printf("%v", *repo.Description)
		}

		fmt.Printf("\n")
	}
}
