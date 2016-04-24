package main

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

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

func option(max int, perPage int) *SearchOpt {
	return &SearchOpt{
		sort:    "best match",
		order:   "desc",
		query:   "ghs",
		perPage: perPage,
		max:     max}
}

type pageTestReulst struct {
	repoLen  int
	lastPage int
	max      int
}

func pageTest(max int, total int, perPage int, lastPage int) *pageTestReulst {
	link := headerLink(perPage, lastPage)
	opts := option(max, perPage)

	url, _ := url.Parse(server.URL)

	repo = NewRepo(NewSearch(opts, url, ""))

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Link", link)
		fmt.Fprintf(w, `{"total_count": %d, "items": [{"id":1},{"id":2}]}`, total)
	})
	repos, lastPage, max := repo.search.First()

	return &pageTestReulst{
		repoLen:  len(repos),
		lastPage: lastPage,
		max:      max,
	}
}

func TestSearch_Page(t *testing.T) {
	result := pageTest(1000, 1000, 100, 3)
	want := &pageTestReulst{
		repoLen:  2,
		lastPage: 3,
		max:      1000,
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("TestSearch_Page returned %+v, want %+v", result, want)
	}
}

// func TestSearch_Link(t *testing.T) {

// 	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		testFormValues(t, r, values{
// 			"q":        "ghs",
// 			"sort":     "best match",
// 			"order":    "desc",
// 			"page":     "1",
// 			"per_page": "100",
// 		})
// 		w.Header().Add("Link", link)
// 		fmt.Fprint(w, `{"total_count": 1000, "items": [{"id":1},{"id":2}]}`)
// 	})

// 	repos, lastPage, max := repo.search.First()

// 	if !(max == 1000 && lastPage == 3 && len(repos) == 2) {
// 		t.Errorf("repo.search.First() max %+v, lastPage %+v, len(repos) %+v\n", max, lastPage, len(repos))
// 	}
// }
