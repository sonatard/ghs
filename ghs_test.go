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
