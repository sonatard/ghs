package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/sona-tar/ghs/debug"
	"os"
)

type GhsOptions struct {
	Sort       string `short:"s"  long:"sort"      description:"The sort field. 'stars', 'forks', or 'updated'." default:"best match"`
	Order      string `short:"o"  long:"order"     description:"The sort order. 'asc' or 'desc'." default:"desc"`
	Language   string `short:"l"  long:"language"  description:"searches repositories based on the language they’re written in."`
	User       string `short:"u"  long:"user"      description:"limits searches to a specific user name."`
	Repository string `short:"r"  long:"repo"      description:"limits searches to a specific repository."`
	Version    bool   `short:"v"  long:"version"   description:"print version infomation and exit."`
}

func GhsOptionParser() ([]string, GhsOptions) {
	var opts GhsOptions
	parser := flags.NewParser(&opts, flags.HelpFlag)

	parser.Name = "ghs"
	parser.Usage = "[OPTION] \"QUERY\""
	args, err := parser.Parse()

	printGhsOption(args, opts)

	if err != nil {
		printGhsHelp(parser)
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("ghs %s\n", version)
		os.Exit(0)
	}

	return args, opts
}

func printGhsOption(args []string, opts GhsOptions) {
	debug.Printf("args = %v\n", args)

	debug.Printf("cmd option sort = %s\n", opts.Sort)
	debug.Printf("cmd option order = %s\n", opts.Order)

	debug.Printf("cmd option language = %s\n", opts.Language)
	debug.Printf("cmd option User = %s\n", opts.User)
	debug.Printf("cmd option Repository = %s\n", opts.Repository)
	debug.Printf("cmd option Version = %s\n", opts.Version)
}

func printGhsHelp(parser *flags.Parser) {
	parser.WriteHelp(os.Stdout)
	fmt.Printf("\n")
	fmt.Printf("Github search APIv3 QUERY infomation:\n")
	fmt.Printf("   https://developer.github.com/v3/search/\n")
	fmt.Printf("   https://help.github.com/articles/searching-repositories/\n")
	fmt.Printf("\n")
	fmt.Printf("Version:\n")
	fmt.Printf("   ghs %s\n", version)
}
