package main

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type pageTestReulst struct {
	printLastPage int
	printMax      int
}

func TestSearch_Page(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}
	// pageTest: max,      total, perPage
	// Result:   lastPage, max
	// normal test
	assert(pageTest(1000, 1000, 100), &pageTestReulst{10, 1000})
	// total test
	assert(pageTest(1000, 100, 100), &pageTestReulst{1, 100})
	assert(pageTest(1000, 101, 100), &pageTestReulst{2, 101})
	assert(pageTest(1000, 500, 100), &pageTestReulst{5, 500})
	assert(pageTest(1000, 1, 100), &pageTestReulst{1, 1})
	assert(pageTest(2000, 1000, 100), &pageTestReulst{10, 1000})
	// max test
	assert(pageTest(100, 1000, 100), &pageTestReulst{1, 100})
	assert(pageTest(101, 1000, 100), &pageTestReulst{2, 101})
	assert(pageTest(500, 1000, 100), &pageTestReulst{5, 500})
	assert(pageTest(1, 1000, 100), &pageTestReulst{1, 1})
	assert(pageTest(2000, 1000, 100), &pageTestReulst{10, 1000})
	// perPage test
	assert(pageTest(1000, 1000, 1), &pageTestReulst{1000, 1000})
	assert(pageTest(1000, 1000, 2), &pageTestReulst{500, 1000})
	assert(pageTest(1000, 1000, 10), &pageTestReulst{100, 1000})
	assert(pageTest(1000, 1000, 50), &pageTestReulst{20, 1000})
	assert(pageTest(1000, 1000, 1000), &pageTestReulst{1, 1000})
}

func pageTest(max int, total int, perPage int) *pageTestReulst {
	Setup()
	defer Teardown()

	// create input
	lastPage := ((total - 1) / perPage) + 1
	url, _ := url.Parse(server.URL)
	repo = NewRepo(NewSearch(option(max, perPage, url, "")))

	// create output
	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Link", headerLink(perPage, lastPage))
		fmt.Fprintf(w, `{"total_count": %d, "items": [{"id":1}]}`, total)
	})

	// test
	_, printLastPage, printMax, _ := repo.search.First()

	return &pageTestReulst{
		printLastPage: printLastPage,
		printMax:      printMax,
	}
}

func headerLink(perPage int, lastPage int) string {
	link := func(per int, page int, rel string) string {
		url := "https://api.github.com/search/repositories"
		query := fmt.Sprintf("?q=ghs&per_page=%d&page=%d", per, page)
		link := "<" + url + query + ">"
		return link + "; " + fmt.Sprintf(`rel="%s"`, rel)
	}
	next := link(perPage, 2, "next")
	last := link(perPage, lastPage, "last")
	return "Link: " + next + ", " + last

}

func option(max int, perPage int, url *url.URL, token string) *SearchOpt {
	return &SearchOpt{
		sort:    "best match",
		order:   "desc",
		query:   "ghs",
		perPage: perPage,
		max:     max,
		baseURL: url,
		token:   token,
	}
}

type requestTestReulst struct {
	repolen  int
	printMax int
}

func TestSearch_Request(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}

	var handler func(w http.ResponseWriter, r *http.Request)
	// Normal response
	handler = func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "ghs",
			"sort":     "best match",
			"order":    "desc",
			"page":     "1",
			"per_page": "100",
		})
		var items []string
		for i := 1; i < 100+1; i++ {
			items = append(items, fmt.Sprintf("{\"id\":%d}", i))
		}
		fmt.Fprintf(w, `{"total_count": 102, "items": [%s]}`, strings.Join(items, ","))
	}
	ret, _ := firstRequestTest(t, handler)
	assert(ret, &requestTestReulst{100, 102})

	// Invalid response
	handler = func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "ghs",
			"sort":     "best match",
			"order":    "desc",
			"page":     "1",
			"per_page": "100",
		})
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
	}
	ret, err := firstRequestTest(t, handler)
	assert(strings.Contains(err.Error(), "404"), true)
	assert(ret, &requestTestReulst{0, 0})

	// Normal response
	handler = func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "ghs",
			"sort":     "best match",
			"order":    "desc",
			"page":     "2",
			"per_page": "100",
		})
		fmt.Fprint(w, `{"total_count": 102, "items": [{"id":1},{"id":2}]}`)
	}
	repoNum, _ := secondRequestTest(t, handler)
	assert(repoNum, 2)
	// Invalid response
	handler = func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "ghs",
			"sort":     "best match",
			"order":    "desc",
			"page":     "2",
			"per_page": "100",
		})
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
	}
	repoNum, err = secondRequestTest(t, handler)
	assert(strings.Contains(err.Error(), "404"), true)
	assert(repoNum, 0)
}

func firstRequestTest(t *testing.T, handler func(http.ResponseWriter, *http.Request)) (*requestTestReulst, error) {
	Setup()
	defer Teardown()

	// create input
	max := 1000
	perPage := 100
	url, _ := url.Parse(server.URL)
	repo = NewRepo(NewSearch(option(max, perPage, url, "")))

	// create output
	mux.HandleFunc("/search/repositories", handler)

	// test
	repos, _, printMax, err := repo.search.First()

	return &requestTestReulst{
		repolen:  len(repos),
		printMax: printMax,
	}, err
}

func secondRequestTest(t *testing.T, handler func(http.ResponseWriter, *http.Request)) (int, error) {
	Setup()
	defer Teardown()

	// create input
	max := 1000
	perPage := 100
	url, _ := url.Parse(server.URL)
	repo = NewRepo(NewSearch(option(max, perPage, url, "abcdefg")))

	// create output
	mux.HandleFunc("/search/repositories", handler)

	// test
	repos, err := repo.search.Exec(2)

	return len(repos), err
}
