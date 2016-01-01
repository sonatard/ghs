package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"net/url"
	"os"
)

type Flags struct {
	args    []string
	cmdOpts CmdOptions
	parser  *flags.Parser
}

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

func (f *Flags) ParseOption() (bool, int, *SearchInfo, *url.URL, string) {
	if f.cmdOpts.Version {
		fmt.Printf("ghs %s\n", Version)
		checkVersion(Version)
		return true, ExitCodeOK, nil, nil, ""
	}

	if (f.cmdOpts.User == "" && f.cmdOpts.Repository == "" && f.cmdOpts.Language == "") && len(f.args) == 0 {
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

	query := f.buildQuery()

	return false, ExitCodeOK,
		&SearchInfo{sort: f.cmdOpts.Sort,
			order:   f.cmdOpts.Order,
			query:   query,
			perPage: 100,
			max:     f.cmdOpts.Max},
		baseURL, f.cmdOpts.Token
}

func (f *Flags) buildQuery() string {
	query := ""

	for _, arg := range f.args {
		query += " " + arg
	}

	if f.cmdOpts.Fields != "" {
		query += " in:" + f.cmdOpts.Fields
	}
	if f.cmdOpts.Language != "" {
		query += " language:" + f.cmdOpts.Language
	}
	if f.cmdOpts.User != "" {
		query += " user:" + f.cmdOpts.User
	}
	if f.cmdOpts.Repository != "" {
		query += " repo:" + f.cmdOpts.Repository
	}

	return query
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
