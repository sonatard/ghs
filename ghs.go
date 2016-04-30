package main

import (
	"log"
	"os"

	"github.com/google/go-github/github"
)

// Version is ghs version number
const Version string = "0.0.7"

const (
	// ExitCodeOK is 0
	ExitCodeOK = iota
	// ExitCodeError is 1
	ExitCodeError
)

func main() {
	flags, err := NewFlags(os.Args[1:])
	if err != nil {
		flags.printHelp()
		os.Exit(ExitCodeError)
	}

	exit, exitCode, sOpt, url, token := flags.ParseOption()
	if exit {
		os.Exit(exitCode)
	}

	repo := NewRepo(NewSearch(sOpt, url, token))
	reposChan, oneRequestFin := repo.Search()

	Debug("main thread select start...\n")
	var repos []github.Repository
	for {
		select {
		case oneReqRepos := <-reposChan:
			Debug("main thread chan reposChan\n")
			Debug("main thread oneReqRepos length %d\n", len(oneReqRepos))

			repos = append(repos, oneReqRepos...)
			Debug("main thread repos length %d\n", len(repos))
		case <-oneRequestFin:
			Debug("main thread chan fin\n")
			end := repo.Print(repos)
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
