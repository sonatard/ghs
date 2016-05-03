package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {

	//	setup()

	code := m.Run()
	//	teardown()

	os.Exit(code)
}

var (
	repo   *Repo
	server *httptest.Server
	mux    *http.ServeMux
)

func Setup() {
	os.Setenv("GHS_PRINT", "no")
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func Teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}
