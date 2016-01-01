package main

import (
	"fmt"
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

	reposBuff, fin := repo.SearchRepository()

	Debug("main thread select start...\n")
	for {
		select {
		case repos := <-reposBuff:
			Debug("main thread chan reposBuff\n")
			Debug("main thread repos length %d\n", len(repos))

			repo.repos = append(repo.repos, repos...)
			Debug("main thread repo.repos length %d\n", len(repo.repos))
		case <-fin:
			Debug("main thread chan fin\n")
			end := repo.PrintRepository()
			if end {
				Debug("over max\n")
				return
			}

			return
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
