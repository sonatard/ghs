package main

import (
	"fmt"
	"os"

	"github.com/tcnksm/go-latest"
)

func CheckVersion(ver string) {
	if os.Getenv("GHS_PRINT") != "no" {
		githubTag := &latest.GithubTag{
			Owner:      "sona-tar",
			Repository: "ghs",
		}
		res, _ := latest.Check(githubTag, ver)
		if res.Outdated {
			fmt.Printf(fmt.Sprintf("%s is not latest, you should upgrade to %s\n", ver, res.Current))
			fmt.Printf("-> $ brew update && brew upgrade sona-tar/tools/ghs")
		}
	}
}
