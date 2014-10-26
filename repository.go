package main

import (
	"fmt"
	"github.com/google/go-github/github"
)

func SearchRepository(sort string, order string, query string) []github.Repository {
	client := github.NewClient(nil)
	searchOpts := &github.SearchOptions{
		Sort:  sort,
		Order: order,
		// TextMatch: true,
		// ListOptions: github.ListOptions{Page: 1, PerPage: 1},
	}

	searchResult, _, err := client.Search.Repositories(query, searchOpts)
	if err != nil {
		fmt.Printf("repository search error")
	}

	return searchResult.Repositories
}

func PrintRepository(repos []github.Repository) {
	for _, repo := range repos {
		if repo.FullName != nil {
			fmt.Printf("%-40v", *repo.FullName)
		}
		if repo.Description != nil {
			fmt.Printf("%v", *repo.Description)
		}

		fmt.Printf("\n")
	}
}
