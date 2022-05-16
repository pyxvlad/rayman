package pacman_test

import (
	"strings"
	"testing"

	"gitlab.com/rayone121/rayman/pacman"
)

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

func TestParseRepositoryFile(t *testing.T) {

	// TODO: add the file in testdata/ instead of assuming an Arch Linux system
	packages, err := pacman.ParseRepositoryFile("core")
	if err != nil {
		t.Fatal("couldn't parse the repository")
	}

	found := false
	for _, p := range packages {
		if p.Name == "linux" && p.Repository == "core" {
			found = true
		}
	}

	if !found {
		t.Fatal("something went wrong... couldn't find 'linux' package...")
	}
}
