package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/google/go-github/github"
)

type Repo struct {
	search     *Search
	maxItem    int
	printCount int
}

func NewRepo(s *Search) *Repo {
	return &Repo{s, 0, 0}
}

func (r *Repo) Search() (<-chan []github.Repository, <-chan bool) {
	var wg sync.WaitGroup
	reposBuff := make(chan []github.Repository, 1000)
	fin := make(chan bool)

	// 1st search
	repos, lastPage, max := r.search.First()
	r.maxItem = max

	// notify main thread of first search result
	reposBuff <- repos

	// 2nd - 10th search
	go func() {
		for page := 2; page < lastPage+1; page++ {
			Debug("sub thread %d\n", page)
			wg.Add(1)
			go func(p int) {
				// notify main thread of 2nd - 10th search result
				reposBuff <- r.search.Exec(p)
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

func (r *Repo) Print(repos []github.Repository) (end bool) {
	Debug("repos length %d\n", len(repos))
	repoNameMaxLen := 0
	for _, repo := range repos {
		repoNamelen := len(*repo.FullName)
		if repoNamelen > repoNameMaxLen {
			repoNameMaxLen = repoNamelen
		}
	}

	printLine := func(repo github.Repository) {
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
	}

	for _, repo := range repos {
		printLine(repo)
		r.printCount++
		Debug("printCount %d, r.maxItem %d\n", r.printCount, r.maxItem)
		if r.printCount >= r.maxItem {
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
