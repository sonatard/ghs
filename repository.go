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
		fmt.Printf("Repository not Found\n")
	}

	return searchResult.Repositories
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
