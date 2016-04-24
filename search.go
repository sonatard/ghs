package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Search struct {
	client *github.Client
	option *SearchOpt
}

type SearchOpt struct {
	sort    string
	order   string
	query   string
	max     int
	perPage int
}

func NewSearch(opt *SearchOpt, baseURL *url.URL, token string) *Search {
	var tc *http.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	cli := github.NewClient(tc)

	if baseURL != nil {
		cli.BaseURL = baseURL
	}

	return &Search{client: cli, option: opt}
}

func repoSearch(client *github.Client, page int, opt *SearchOpt) (*github.RepositoriesSearchResult, *github.Response, error) {
	opts := &github.SearchOptions{
		Sort:        opt.sort,
		Order:       opt.order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: opt.perPage, Page: page},
	}
	return client.Search.Repositories(opt.query, opts)
}

func (s *Search) First() (repos []github.Repository, lastPage int, maxItem int) {
	ret, resp, err := repoSearch(s.client, 1, s.option)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", s.option.query)
		fmt.Println(err)
		os.Exit(1)
	}

	Debug("main thread repos length %d\n", len(ret.Repositories))

	max := s.option.max
	Debug("Total = %d\n", *ret.Total)
	Debug("s.option.max = %d\n", s.option.max)
	if *ret.Total < max {
		max = *ret.Total
	}

	last := ((max - 1) / s.option.perPage) + 1
	Debug("resp.LastPage = %d\n", resp.LastPage)
	Debug("last = %d\n", last)
	if resp.LastPage < last {
		last = resp.LastPage
	}

	return ret.Repositories, last, max
}

func (s *Search) Exec(page int) (repos []github.Repository) {
	Debug("Page%d go func search start\n", page)

	Debug("Page%d query : %s\n", page, s.option.query)
	ret, _, err := repoSearch(s.client, page, s.option)
	if err != nil {
		fmt.Printf("Search Error!! query : %s\n", s.option.query)
		fmt.Println(err)
	}
	Debug("Page%d go func search result length %d\n", page, len(ret.Repositories))
	Debug("Page%d go func search end\n", page)

	return ret.Repositories
}
