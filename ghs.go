package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"os"
)

const Version string = "0.0.7"

const (
	ExitCodeOK = iota
	ExitCodeError
)

func main() {
	flags, err := NewFlags()
	if err != nil {
		flags.printHelp()
		os.Exit(ExitCodeError)
	}

	exit, exitCode, searchInfo, url, token := flags.ParseOption()
	if exit {
		os.Exit(exitCode)
	}

	repo, err := NewRepo(searchInfo, url, token)
	if err != nil {
		fmt.Printf("Option error\n")
	}

	reposBuff, one_request_fin := repo.SearchRepository()

	Debug("main thread select start...\n")
	var repos []github.Repository
	for {
		select {
		case one_req_repos := <-reposBuff:
			Debug("main thread chan reposBuff\n")
			Debug("main thread one_req_repos length %d\n", len(one_req_repos))

			repos = append(repos, one_req_repos...)
			Debug("main thread repos length %d\n", len(repos))
		case <-one_request_fin:
			Debug("main thread chan fin\n")
			end := repo.PrintRepository(repos)
			if end {
				Debug("over max\n")
				return
			}
		}
	}
}

// Debug display values when DEBUG mode
// This is used only for developer
func Debug(format string, args ...interface{}) {
	if os.Getenv("GHS_DEBUG") != "" {
		log.Printf(format, args...)
	}
}
