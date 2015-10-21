package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"net/url"
)

func SearchRepository(sort string, order string, max int, enterprise string, query string) []github.Repository {
	client := github.NewClient(nil)

	if enterprise != "" {
		baseURL, err := url.Parse(enterprise)
		if err == nil {
			client.BaseURL = baseURL
		} else {
			fmt.Printf("%s cannot parse\n", enterprise)
		}
	}

	perPage := 100

	if max < 100 {
		perPage = max
	}

	searchOpts := &github.SearchOptions{
		Sort:        sort,
		Order:       order,
		TextMatch:   false,
		ListOptions: github.ListOptions{PerPage: perPage},
	}

	var allRepos []github.Repository
	i := 0

	for {
		searchResult, resp, err := client.Search.Repositories(query, searchOpts)
		if err != nil {
			fmt.Printf("Repository not Found\n")
		}

		i++
		allRepos = append(allRepos, searchResult.Repositories...)

		if resp.NextPage == 0 || (i*perPage) >= max {
			break
		}

		searchOpts.ListOptions.Page = resp.NextPage
	}

	return allRepos
}

func PrintRepository(repos []github.Repository) {
	repoNameMaxLen := 0
	for _, repo := range repos {
		repoNamelen := len(*repo.FullName)
		if repoNamelen > repoNameMaxLen {
			repoNameMaxLen = repoNamelen
		}
	}
	for _, repo := range repos {
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
