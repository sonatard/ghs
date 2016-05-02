package main

import (
	"net/http"
	"net/url"

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
	baseURL *url.URL
	token   string
}

func NewSearch(opt *SearchOpt) *Search {
	var tc *http.Client

	if opt.token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: opt.token})
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}

	cli := github.NewClient(tc)

	if opt.baseURL != nil {
		cli.BaseURL = opt.baseURL
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

func (s *Search) First() (repos []github.Repository, lastPage int, maxItem int, err error) {
	ret, resp, err := repoSearch(s.client, 1, s.option)
	if err != nil {
		return nil, 0, 0, err
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

	return ret.Repositories, last, max, nil
}

func (s *Search) Exec(page int) (repos []github.Repository, err error) {
	Debug("Page%d go func search start\n", page)

	Debug("Page%d query : %s\n", page, s.option.query)
	ret, _, err := repoSearch(s.client, page, s.option)
	if err != nil {
		return nil, err
	}
	Debug("Page%d go func search result length %d\n", page, len(ret.Repositories))
	Debug("Page%d go func search end\n", page)

	return ret.Repositories, nil
}
