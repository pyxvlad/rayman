package aurweb_test

import (
	"fmt"
	"testing"

	"gitlab.com/rayone121/rayman/aurweb"
)

// NOTE: right now it just checks if it can do a request
func TestInfo(t *testing.T) {
	s := []string{"bar", "neovim-git"}
	info, err := aurweb.Info(s)
	if err != nil {
		t.Fail()

	}
	fmt.Printf("%#v", info)
}

// NOTE: right now it just checks if it can do a request
func TestSearch(t *testing.T) {
	search, err := aurweb.Search("name-desc", "gitlab-cli")
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%#v", search)
}
