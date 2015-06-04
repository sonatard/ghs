package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/tcnksm/go-latest"
)

func checkVersion(ver string) {
	githubTag := &latest.GithubTag{
		Owner:      "sona-tar",
		Repository: "ghs",
	}
	res, _ := latest.Check(githubTag, ver)
	if res.Outdated {
		redstr := ansi.ColorFunc("red")
		fmt.Printf(redstr(fmt.Sprintf("%s is not latest, you should upgrade to %s\n", Version, res.Current)))
		fmt.Printf(redstr("-> $ brew update && brew upgrade sona-tar/tools/ghs"))
	}
}
