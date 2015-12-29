package main

import (
	"fmt"
	"log"
	"os"
)

const Version string = "0.0.6"

func buildQuery(args []string, opts GhsOptions) string {
	query := ""

	for _, arg := range args {
		query += arg
	}

	if opts.Fields != "" {
		query += " in:" + opts.Fields
	}
	if opts.Language != "" {
		query += " language:" + opts.Language
	}
	if opts.User != "" {
		query += " user:" + opts.User
	}
	if opts.Repository != "" {
		query += " repo:" + opts.Repository
	}

	return query
}

func main() {
	args, opts := GhsOptionParser()
	query := buildQuery(args, opts)

	repo, err := NewRepo(opts.Sort, opts.Order, opts.Max, opts.Enterprise, opts.Token, query)
	if err != nil {
		fmt.Printf("Option error\n")
	}

	reposBuff, fin := repo.SearchRepository()

	for {
		select {
		case repo.repos = <-reposBuff:
			Debug("print\n")
			end := repo.PrintRepository()
			if end {
				Debug("over max\n")
				return
			}
		case <-fin:
			Debug("fin\n")
			return
		}
	}

}

// Debug display values when DEBUG mode
// This is used only for developer
func Debug(v ...interface{}) {
	if os.Getenv("GHS_DEBUG") != "" {
		log.Println(v...)
	}
}
