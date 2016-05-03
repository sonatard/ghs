package main

import (
	"fmt"
	"log"
	"os"

	"github.com/motemen/go-gitconfig"
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
	num, err := ghs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(ExitCodeError)
	}
	Debug("Print %d\n", num)
}

func ghs(args []string) (int, error) {
	flags, err := NewFlags(args)
	// --help or error
	if err != nil {
		Debug("Error : help or parse error\n")
		flags.PrintHelp()
		CheckVersion(Version)
		os.Exit(ExitCodeError)
	}

	version, exitCode, sOpt := flags.ParseOption()
	// --version
	if version {
		fmt.Printf("ghs %s\n", Version)
		CheckVersion(Version)
		os.Exit(ExitCodeOK)

	}
	// error options
	if exitCode == ExitCodeError {
		Debug("Error : Parse option error flags.ParseOption()\n")
		flags.PrintHelp()
		CheckVersion(Version)
		os.Exit(ExitCodeError)
	}
	getToken := func(optsToken string) string {
		// -t or --token option
		if optsToken != "" {
			Debug("Github token get from option value\n")
			return optsToken
		}

		// GITHUB_TOKEN environment
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			Debug("Github token get from environment value\n")
			return token
		}

		// github.token in gitconfig
		if token, err := gitconfig.GetString("github.token"); err == nil {
			Debug("Github token get from gitconfig value\n")
			return token
		}

		Debug("Github token not found\n")
		return ""
	}
	sOpt.token = getToken(sOpt.token)

	repo := NewRepo(NewSearch(sOpt))
	reposChan, errChan := repo.Search()

	Debug("main thread select start...\n")
	reposNum := 0

	for {
		select {
		case oneReqRepos := <-reposChan:
			Debug("main thread chan reposChan\n")
			Debug("main thread oneReqRepos length %d\n", len(oneReqRepos))
			var end bool
			end, reposNum = repo.Print(oneReqRepos)
			Debug("reposNum : %d\n", reposNum)
			if end {
				Debug("over max\n")
				return reposNum, nil

			}
		case err := <-errChan:
			Debug("main thread chan err\n")
			return 0, err
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
