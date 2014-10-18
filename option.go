package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/sona-tar/ghs/debug"
	"os"
)

type GhsOptions struct {
	Sort  string `short:"s" long:"sort"  description:" The sort field. 'stars', 'forks', or 'updated'." default:"best match"`
	Order string `short:"o" long:"order" description:" The sort order. 'asc' or 'desc'." default:"desc"`
}

func GhsOptionParser() (string, string, string) {
	var opts GhsOptions
	parser := flags.NewParser(&opts, flags.HelpFlag)

	parser.Name = "ghs"
	parser.Usage = "[OPTION] \"QUERY\"(The search keywords, as well as any qualifiers.)"
	args, err := parser.Parse()

	if len(args) != 1 || err != nil {
		parser.WriteHelp(os.Stdout)
		fmt.Printf("\n")
		fmt.Printf("Github search APIv3 QUERY infomation:\n")
		fmt.Printf("   https://developer.github.com/v3/search/\n")
		fmt.Printf("   https://help.github.com/articles/searching-repositories/\n")
		os.Exit(1)
	}

	query := args[0]
	debug.Printf("cmd option sort = %s\n", opts.Sort)
	debug.Printf("cmd option order = %s\n", opts.Order)
	debug.Printf("cmd args query = %s\n", query)

	return opts.Sort, opts.Order, query
}
