package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

// Flags is args, opts and parser
type Flags struct {
	args    []string
	cmdOpts CmdOptions
	parser  *flags.Parser
}

// CmdOptions is ghs command line option list
type CmdOptions struct {
	Fields     string `short:"f"  long:"fields"     description:"limits what fields are searched. 'name', 'description', or 'readme'."`
	Fork       string `short:"k"  long:"fork"       description:"Forked repositories icluded in results. 'true', 'only' or 'false'."`
	Sort       string `short:"s"  long:"sort"       description:"The sort field. 'stars', 'forks', or 'updated'." default:"best match"`
	Order      string `short:"o"  long:"order"      description:"The sort order. 'asc' or 'desc'." default:"desc"`
	Language   string `short:"l"  long:"language"   description:"searches repositories based on the language they’re written in."`
	User       string `short:"u"  long:"user"       description:"limits searches to a specific user name."`
	Repository string `short:"r"  long:"repo"       description:"limits searches to a specific repository."`
	Max        int    `short:"m"  long:"max"        description:"limits number of result. range 1-1000" default:"100"`
	Version    bool   `short:"v"  long:"version"    description:"print version infomation and exit."`
	Enterprise string `short:"e"  long:"enterprise" description:"search from github enterprise."`
	Token      string `short:"t"  long:"token"      description:"Github API token to avoid Github API rate "`
	Help       string `short:"h"  long:"help"       description:"Show this help message"`
}

func NewFlags(args []string) (*Flags, error) {
	var opts CmdOptions
	parser := flags.NewParser(&opts, flags.None)
	parser.Name = "ghs"
	parser.Usage = "[OPTION] \"QUERY\""
	args, err := parser.ParseArgs(args)
	return &Flags{args, opts, parser}, err
}

func (f *Flags) ParseOption() (bool, int, *SearchOpt) {
	if f.cmdOpts.Version {
		return true, ExitCodeOK, nil
	}

	errorQuery := func(opts CmdOptions, args []string) bool {
		noargs := len(args) == 0
		noopt := opts.User == "" && opts.Repository == "" && opts.Language == ""
		return noargs && noopt
	}

	if errorQuery(f.cmdOpts, f.args) {
		Debug("Error : noargs & noopt\n")
		return false, ExitCodeError, nil
	}

	if f.cmdOpts.Max < 1 || f.cmdOpts.Max > 1000 {
		return false, ExitCodeError, nil
	}

	var baseURL *url.URL
	if f.cmdOpts.Enterprise != "" {
		eURL, err := url.Parse(f.cmdOpts.Enterprise)
		Debug("%#v\n", eURL)
		if err != nil {
			Debug(`Error : Parse "%v"`+"\n", f.cmdOpts.Enterprise)
			return false, ExitCodeError, nil
		}
		baseURL = eURL
	}

	buildQuery := func(opts CmdOptions, args []string) string {
		var query []string

		if opts.Fields != "" {
			query = append(query, "in:"+opts.Fields)
		}
		if opts.Fork != "" {
			query = append(query, "fork:"+opts.Fork)
		}
		if opts.Language != "" {
			query = append(query, "language:"+opts.Language)
		}
		if opts.User != "" {
			query = append(query, "user:"+opts.User)
		}
		if opts.Repository != "" {
			query = append(query, "repo:"+opts.Repository)
		}
		query = append(query, args...)
		return strings.Join(query, " ")
	}

	return false, ExitCodeOK,
		&SearchOpt{
			sort:    f.cmdOpts.Sort,
			order:   f.cmdOpts.Order,
			query:   buildQuery(f.cmdOpts, f.args),
			perPage: 100,
			max:     f.cmdOpts.Max,
			baseURL: baseURL,
			token:   f.cmdOpts.Token,
		}
}

func (f *Flags) PrintHelp() {
	if os.Getenv("GHS_PRINT") != "no" {
		f.parser.WriteHelp(os.Stdout)
		fmt.Printf("\n")
		fmt.Printf("Github search APIv3 QUERY infomation:\n")
		fmt.Printf("  https://developer.github.com/v3/search/\n")
		fmt.Printf("  https://help.github.com/articles/searching-repositories/\n")
		fmt.Printf("\n")
		fmt.Printf("Version:\n")
		fmt.Printf("  ghs %s (https://github.com/sona-tar/ghs.git)\n", Version)
		fmt.Printf("\n")
	}
}
