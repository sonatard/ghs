package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/motemen/go-gitconfig"
)

// Flags is args, opts and parser
type Flags struct {
	args    []string
	cmdOpts CmdOptions
	parser  *flags.Parser
}

// CmdOptions is ghs command line option list
type CmdOptions struct {
	Fields     string `short:"f"  long:"fields"     description:"limits what fields are searched. 'name', 'description', or 'readme'." default:"name,description"`
	Sort       string `short:"s"  long:"sort"       description:"The sort field. 'stars', 'forks', or 'updated'." default:"best match"`
	Order      string `short:"o"  long:"order"      description:"The sort order. 'asc' or 'desc'." default:"desc"`
	Language   string `short:"l"  long:"language"   description:"searches repositories based on the language they’re written in."`
	User       string `short:"u"  long:"user"       description:"limits searches to a specific user name."`
	Repository string `short:"r"  long:"repo"       description:"limits searches to a specific repository."`
	Max        int    `short:"m"  long:"max"        description:"limits number of result. range 1-1000" default:"100"`
	Version    bool   `short:"v"  long:"version"    description:"print version infomation and exit."`
	Enterprise string `short:"e"  long:"enterprise" description:"search from github enterprise."`
	Token      string `short:"t"  long:"token"      description:"Github API token to avoid Github API rate "`
}

func NewFlags() (*Flags, error) {
	var opts CmdOptions
	parser := flags.NewParser(&opts, flags.HelpFlag)

	parser.Name = "ghs"
	parser.Usage = "[OPTION] \"QUERY\""
	args, err := parser.Parse()

	return &Flags{args, opts, parser}, err
}

func (f *Flags) ParseOption() (bool, int, *SearchOpt, *url.URL, string) {
	if f.cmdOpts.Version {
		fmt.Printf("ghs %s\n", Version)
		checkVersion(Version)
		return true, ExitCodeOK, nil, nil, ""
	}

	errorQuery := func(opts CmdOptions, args []string) bool {
		noargs := len(args) == 0
		noopt := opts.User == "" && opts.Repository == "" && opts.Language == ""
		return noargs && noopt
	}

	if errorQuery(f.cmdOpts, f.args) {
		f.printHelp()
		checkVersion(Version)
		return true, ExitCodeError, nil, nil, ""
	}

	if f.cmdOpts.Max < 1 || f.cmdOpts.Max > 1000 {
		f.printHelp()
		checkVersion(Version)
		return true, ExitCodeError, nil, nil, ""
	}

	var baseURL *url.URL
	if f.cmdOpts.Enterprise != "" {
		url, err := url.Parse(f.cmdOpts.Enterprise)
		if err != nil {
			f.printHelp()
			checkVersion(Version)
			return true, ExitCodeError, nil, nil, ""
		}
		baseURL = url
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

	buildQuery := func(opts CmdOptions, args []string) string {
		query := ""

		for _, arg := range args {
			query += " " + arg
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

	return false, ExitCodeOK,
		&SearchOpt{
			sort:    f.cmdOpts.Sort,
			order:   f.cmdOpts.Order,
			query:   buildQuery(f.cmdOpts, f.args),
			perPage: 100,
			max:     f.cmdOpts.Max},
		baseURL, getToken(f.cmdOpts.Token)
}

func (f *Flags) printHelp() {
	f.parser.WriteHelp(os.Stdout)
	fmt.Printf("\n")
	fmt.Printf("Github search APIv3 QUERY infomation:\n")
	fmt.Printf("   https://developer.github.com/v3/search/\n")
	fmt.Printf("   https://help.github.com/articles/searching-repositories/\n")
	fmt.Printf("\n")
	fmt.Printf("Version:\n")
	fmt.Printf("   ghs %s (https://github.com/sona-tar/ghs.git)\n", Version)
	fmt.Printf("\n")
	checkVersion(Version)
}
