package aurweb_test

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"gitlab.com/rayone121/rayman/aurweb"
)

func TestInfo_HttpErr(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.Deactivate)

	val := url.Values{}
	val.Add("type", "info")
	val.Add("arg[]", "neovim-git")
	val.Add("v", "5")
	expected_err := errors.New("http_error")
	httpmock.RegisterResponderWithQuery("GET", aurweb.AUR, val.Encode(), httpmock.NewErrorResponder(expected_err))

	_, err := aurweb.Info([]string{"neovim-git"})
	if err.(*url.Error).Err != expected_err {
		t.Fatalf("expected err %#v got %#v", expected_err, err)
	}
}

func TestInfo(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.Deactivate)

	pkg := aurweb.InfoResult{Result: aurweb.Result{Name: "neovim-git"}}

	info := aurweb.InfoResponse{ResultCount: 1, Results: []aurweb.InfoResult{pkg}}
	val := url.Values{}
	val.Add("type", "info")
	val.Add("arg[]", "neovim-git")
	val.Add("v", "5")
	httpmock.RegisterResponderWithQuery("GET", aurweb.AUR, val.Encode(), httpmock.NewJsonResponderOrPanic(200, &info))

	got, err := aurweb.Info([]string{"neovim-git"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(info, got) {
		t.Fatalf("expected %#v got %#v", info, got)
	}
}

func TestSearch_HttpErr(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.Deactivate)

	sortBy := "name-desc"
	pkgName := "rayman"

	val := url.Values{}
	val.Add("type", "search")
	val.Add("by", sortBy)
	val.Add("arg", pkgName)
	val.Add("v", "5")

	expected_err := errors.New("http_error")
	httpmock.RegisterResponderWithQuery("GET", aurweb.AUR, val.Encode(), httpmock.NewErrorResponder(expected_err))

	_, err := aurweb.Search(sortBy, pkgName)
	if err.(*url.Error).Err != expected_err {
		t.Fatalf("expected err %#v got %#v", expected_err, err)
	}
}

func TestSearch_InvalidJson(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.Deactivate)

	sortBy := "name-desc"
	pkgName := "rayman"

	val := url.Values{}
	val.Add("type", "search")
	val.Add("by", sortBy)
	val.Add("arg", pkgName)
	val.Add("v", "5")

	httpmock.RegisterResponderWithQuery("GET", aurweb.AUR, val.Encode(), httpmock.NewStringResponder(200, "invalid json"))

	_, err := aurweb.Search(sortBy, pkgName)
	if syntaxErr := err.(*json.SyntaxError); syntaxErr == nil {
		t.Fatalf("expected json.SyntaError, got %#v", err)
	}
}


func TestSearch_TypeError(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.Deactivate)

	sortBy := "name-desc"
	pkgName := "rayman"

	errText := "type error for the test"
	expected := aurweb.SearchResponse{Type: "error", Error: errText}

	val := url.Values{}
	val.Add("type", "search")
	val.Add("by", sortBy)
	val.Add("arg", pkgName)
	val.Add("v", "5")

	httpmock.RegisterResponderWithQuery("GET", aurweb.AUR, val.Encode(), httpmock.NewJsonResponderOrPanic(200, &expected))

	_, err := aurweb.Search(sortBy, pkgName)
	if err.Error() != errText {
		t.Fatalf("expected err.Error() = %#v got %#v", errText, err.Error())
	}
}

func TestSearch(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.Deactivate)

	sortBy := "name-desc"
	pkgName := "rayman"

	pkg := aurweb.Result{Name: pkgName}
	expected := aurweb.SearchResponse{ResultCount: 1, Results: []aurweb.Result{pkg}, Type: "search"}

	val := url.Values{}
	val.Add("type", "search")
	val.Add("by", sortBy)
	val.Add("arg", pkgName)
	val.Add("v", "5")

	httpmock.RegisterResponderWithQuery("GET", aurweb.AUR, val.Encode(), httpmock.NewJsonResponderOrPanic(200, &expected))

	search, err := aurweb.Search(sortBy, pkgName)
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(search, expected) {
		t.Fatalf("expected %#v got %#v", expected, search)
	}
}
