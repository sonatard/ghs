package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type parseTestReulst struct {
	version  bool
	exitCode int
	sOpt     *SearchOpt
}

func ExampleHelp() {
	_ = testParse("ghs -h")
	// Output:
	// Usage:
	//         ghs [OPTION] "QUERY"

	//         Application Options:
	//         -f, --fields=     limits what fields are searched. 'name', 'description', or 'readme'.
	//                 -s, --sort=       The sort field. 'stars', 'forks', or 'updated'. (default: best match)
	//         -o, --order=      The sort order. 'asc' or 'desc'. (default: desc)
	//         -l, --language=   searches repositories based on the language they’re written in.
	//                 -u, --user=       limits searches to a specific user name.
	//                 -r, --repo=       limits searches to a specific repository.
	//                 -m, --max=        limits number of result. range 1-1000 (default: 100)
	//         -v, --version     print version infomation and exit.
	//                 -e, --enterprise= search from github enterprise.
	//                 -t, --token=      Github API token to avoid Github API rate
	//         -h, --help=       Show this help message

	//         Github search APIv3 QUERY infomation:
	// https://developer.github.com/v3/search/
	// https://help.github.com/articles/searching-repositories/

	// Version:
	//         ghs 0.0.7 (https://github.com/sona-tar/ghs.git)
}

func TestOption_Parse(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}

	// want :exit, exitCode , sOpt, url, token
	assert(testParse("ghs -v"), &parseTestReulst{true, ExitCodeOK, nil})
	assert(testParse("ghs -h"), &parseTestReulst{false, ExitCodeError, nil})

	defaultOpt := SearchOpt{
		sort:    "best match",
		order:   "desc",
		query:   "SEARCH_WORD",
		max:     100,
		perPage: 100,
		baseURL: nil,
		token:   "",
	}

	// normal query test
	wantOpt := defaultOpt
	assert(testParse("ghs SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.sort = "stars"
	assert(testParse("ghs -s stars SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.order = "asc"
	assert(testParse("ghs -o asc SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.max = 1000
	assert(testParse("ghs -m 1000 SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.baseURL, _ = url.Parse("http://test.exmaple/")
	assert(testParse("ghs -e http://test.exmaple/ SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.token = "abcdefg"
	assert(testParse("ghs -t abcdefg SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "in:name SEARCH_WORD"
	assert(testParse("ghs -f name SEARCH_WORD"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "user:sona-tar"
	assert(testParse("ghs -u sona-tar"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "repo:sona-tar/ghs"
	assert(testParse("ghs -r sona-tar/ghs"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "language:golang"
	assert(testParse("ghs -l golang"), &parseTestReulst{false, ExitCodeOK, &wantOpt})

	// no args test
	assert(testParse("ghs"), &parseTestReulst{false, ExitCodeError, nil})
	assert(testParse("ghs -o asc"), &parseTestReulst{false, ExitCodeError, nil})
	wantOpt = defaultOpt
	wantOpt.query = "user:sona-tar"
	assert(testParse("ghs -u sona-tar"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "repo:sona-tar/ghs"
	assert(testParse("ghs -r sona-tar/ghs"), &parseTestReulst{false, ExitCodeOK, &wantOpt})
	wantOpt = defaultOpt
	wantOpt.query = "language:golang"
	assert(testParse("ghs -l golang"), &parseTestReulst{false, ExitCodeOK, &wantOpt})

	// invalid option value test
	assert(testParse("ghs -m 1001 SEARCH_WORD"), &parseTestReulst{false, ExitCodeError, nil})
	assert(testParse("ghs -m 0 SEARCH_WORD"), &parseTestReulst{false, ExitCodeError, nil})
	assert(testParse("ghs -e : SEARCH_WORD"), &parseTestReulst{false, ExitCodeError, nil})
}

func testParse(args_string string) *parseTestReulst {
	args := strings.Split(args_string, " ")[1:]
	flags, _ := NewFlags(args)
	version, exitCode, sOpt := flags.ParseOption()
	return &parseTestReulst{version, exitCode, sOpt}
}
