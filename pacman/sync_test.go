package pacman_test

import (
	"errors"
	"strings"
	"testing"

	"gitlab.com/rayone121/rayman/pacman"
)

type ErrorReader struct{}

func (e ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("kaboom")
}

func TestFromDescWithErrorReader(t *testing.T) {
	_, err := pacman.FromDescReader(ErrorReader{})
	if err == nil {
		t.Fatal("parsed pkg from errorneous reader without returning error")
	}
}

func TestFromDesc(t *testing.T) {
	desc :=
		`
	%NAME%
	rayman
	%VERSION%
	0.1.2
	%DESC%
	Web UI and Go powered AUR helper
	`

	pkg, err := pacman.FromDescReader(strings.NewReader(desc))
	if err != nil {
		t.Fatal("couldn't parse from desc")
	}
	if pkg.Name != "rayman" || pkg.Version != "0.1.2" {
		t.Fatal("pkg has been parsed wrong")
	}
}

func TestParseInexistentSystemRepositoryFile(t *testing.T) {
	_, err := pacman.ParseSystemRepositoryFile("invalid")
	if err == nil {
		t.Fatal("parsed invalid system repository file")
	}
}

func TestParseSystemRepositoryFile(t *testing.T) {

	// TODO: add the file in testdata/ instead of assuming an Arch Linux system
	packages, err := pacman.ParseSystemRepositoryFile("core")
	if err != nil {
		t.Fatal("couldn't parse the repository")
	}

	found := false
	for _, p := range packages {
		if p.Name == "linux" && p.Repository == "core" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("something went wrong... couldn't find 'linux' package...")
	}
}
