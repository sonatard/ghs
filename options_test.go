package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type parseTestReulst struct {
	exit     bool
	exitCode int
	sOpt     *SearchOpt
	url      *url.URL
	token    string
}

func ExampleOption_Parse() {
	// args, exit, exitCode , sOpt, url, token
	_ = testParse("ghs -v")
	// Output:
	// ghs 0.0.7
}

func TestOption_Parse(t *testing.T) {
	assert := func(result interface{}, want interface{}) {
		if want == nil && result != nil {
			t.Errorf("Returned %+v, want nil", result)
		}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Returned %+v, want %+v", result, want)
		}
	}

	// args, exit, exitCode , sOpt, url, token
	result := testParse("ghs -v")
	want := &parseTestReulst{true, 0, nil, nil, ""}
	assert(result, want)
}

func testParse(args_string string) *parseTestReulst {
	args := strings.Split(args_string, " ")[1:]
	flags, _ := NewFlags(args)
	exit, exitCode, sOpt, url, token := flags.ParseOption()
	return &parseTestReulst{exit, exitCode, sOpt, url, token}
}
