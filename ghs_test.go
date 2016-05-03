package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func Test_Ghs(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}

	Setup()
	defer Teardown()

	// Normal response
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Link", HeaderLink(100, 10))
		var items []string
		for i := 1; i < 100+1; i++ {
			items = append(items, fmt.Sprintf(`{"id":%d, "full_name": "test/search_word%d"}`, i, i))
		}
		fmt.Fprintf(w, `{"total_count": 1000, "items": [%s]}`, strings.Join(items, ","))
	}
	mux.HandleFunc("/search/repositories", handler)

	num, err := ghs(strings.Split(fmt.Sprintf("-e %s -m 1000 SEARCH_WORD", server.URL), " "))
	assert(num, 1000)
	assert(err, nil)

	num, err = ghs(strings.Split(fmt.Sprintf("-e %s -m 110 SEARCH_WORD", server.URL), " "))
	assert(num, 110)
	assert(err, nil)

	num, err = ghs(strings.Split("-v", " "))
	assert(num, 0)
	assert(err, nil)

	num, err = ghs(strings.Split("-h", " "))
	assert(num, 0)
	assert(err, errors.New("help or parse error"))

	num, err = ghs(strings.Split("-wrong_option", " "))
	assert(num, 0)
	assert(err, errors.New("help or parse error"))

	num, err = ghs(strings.Split("-s stars", " "))
	assert(num, 0)
	assert(err, errors.New("Parse option error."))
}

func Test_GhsTokenTest(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}

	Setup()
	defer Teardown()

	// Normal response
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Link", HeaderLink(100, 10))
		var items []string
		for i := 1; i < 100+1; i++ {
			items = append(items, fmt.Sprintf(`{"id":%d, "full_name": "test/search_word%d"}`, i, i))
		}
		fmt.Fprintf(w, `{"total_count": 1000, "items": [%s]}`, strings.Join(items, ","))
	}
	mux.HandleFunc("/search/repositories", handler)

	num, err := ghs(strings.Split(fmt.Sprintf("-t abcdefg -e %s -m 100 SEARCH_WORD", server.URL), " "))
	assert(num, 100)
	assert(err, nil)

	os.Setenv("GITHUB_TOKEN", "abcdefg")
	num, err = ghs(strings.Split(fmt.Sprintf("-e %s -m 100 SEARCH_WORD", server.URL), " "))
	assert(num, 100)
	assert(err, nil)
	os.Setenv("GITHUB_TOKEN", "")

	panicIf := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	must := panicIf

	run := func(cmd string, args ...string) error {
		return exec.Command(cmd, args...).Run()
	}

	tmpHome, err := ioutil.TempDir("", "go-gitconfig")
	panicIf(err)

	repoDir := filepath.Join(tmpHome, "repo")
	must(os.Setenv("HOME", tmpHome))
	must(os.MkdirAll(repoDir, 0777))
	must(os.Chdir(repoDir))

	must(run("git", "init"))
	must(run("git", "config", "--global", "github.token", "abcdefg"))

	num, err = ghs(strings.Split(fmt.Sprintf("-e %s -m 100 SEARCH_WORD", server.URL), " "))
	assert(num, 100)
	assert(err, nil)
	must(os.RemoveAll(tmpHome))
}
func Test_GhsInvalidResponse(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}

	Setup()
	defer Teardown()

	// Invalid response
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
	}
	mux.HandleFunc("/search/repositories", handler)

	num, err := ghs(strings.Split(fmt.Sprintf("-e %s -m 100 SEARCH_WORD", server.URL), " "))
	assert(num, 0)
	assert(strings.Contains(err.Error(), "404"), true)
}

func Example_CheckVersion() {
	os.Setenv("GHS_PRINT", "yes")
	defer os.Setenv("GHS_PRINT", "no")

	CheckVersion("0.0.1")
	// Output:
	// 0.0.1 is not latest, you should upgrade to 0.0.7
	// -> $ brew update && brew upgrade sona-tar/tools/ghs
}

func Example_PrintHelp() {
	os.Setenv("GHS_PRINT", "yes")
	defer os.Setenv("GHS_PRINT", "no")
	args := strings.Split("ghs -h", " ")[1:]
	flags, _ := NewFlags(args)
	flags.PrintHelp()
	// Output:
	// Usage:
	//   ghs [OPTION] "QUERY"
	//
	// Application Options:
	//   -f, --fields=     limits what fields are searched. 'name', 'description', or
	//                     'readme'.
	//   -s, --sort=       The sort field. 'stars', 'forks', or 'updated'. (default:
	//                     best match)
	//   -o, --order=      The sort order. 'asc' or 'desc'. (default: desc)
	//   -l, --language=   searches repositories based on the language they’re
	//                     written in.
	//   -u, --user=       limits searches to a specific user name.
	//   -r, --repo=       limits searches to a specific repository.
	//   -m, --max=        limits number of result. range 1-1000 (default: 100)
	//   -v, --version     print version infomation and exit.
	//   -e, --enterprise= search from github enterprise.
	//   -t, --token=      Github API token to avoid Github API rate
	//   -h, --help=       Show this help message
	//
	// Github search APIv3 QUERY infomation:
	//   https://developer.github.com/v3/search/
	//   https://help.github.com/articles/searching-repositories/
	//
	// Version:
	//   ghs 0.0.7 (https://github.com/sona-tar/ghs.git)
}

func Example_PrintfDebug() {
	os.Setenv("GHS_PRINT", "yes")
	Printf("test")
	// Output:
	// test
	os.Setenv("GHS_PRINT", "no")

	os.Setenv("GHS_DEBUG", "yes")

	Debug("test")
	// Output:
	// test
	os.Setenv("GHS_DEBUG", "")
}
