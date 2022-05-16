package pacman

import (
	"bufio"
	"io"
	"os"
	"strings"

	"gitlab.com/rayone121/rayman/util"
)

const descFileFields = 16

type Parser func(io.Reader) (Package, error)

// FromDescReader returns a Package filled with the details parsed from the io.Reader passed to it
func FromDescReader(r io.Reader) (Package, error) {
	reader := bufio.NewReader(r)
	var pkg Package

	lines := make([]string, 0, descFileFields*2)
	l, _, err := reader.ReadLine()
	for err != io.EOF {
		if err != nil {
			return Package{}, err
		}
		line := string(l)

		// make sure the line is not empty
		if len(line) > 1 {
			lines = append(lines, strings.TrimSpace(line))
		}
		l, _, err = reader.ReadLine()
	}

	header := ""

	// TODO: add more fields maybe?
	for _, l := range lines {
		if strings.HasPrefix(l, "%") {
			header = l
			continue
		}

		switch header {
		case "%VERSION%":
			pkg.Version = l
		case "%NAME%":
			pkg.Name = l
		case "%DESC%":
			pkg.Description = l
		}
	}

	return pkg, nil
}

func ParseSystemRepositoryFile(repoName string, parser Parser) ([]Package, error) {
	pkg, err := ParseRepositoryFile("/var/lib/pacman/sync/" + repoName + ".db", parser)
	if err != nil {
		return nil, err
	}
	for i := range pkg {
		pkg[i].Repository = repoName
	}
	return pkg, nil
}

func ParseRepositoryFile(path string, parser Parser) ([]Package, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	packages, err := ParseRepository(file, parser)
	if err != nil {
		util.AssertNoError(file.Close())
		return nil, err
	}

	return packages, file.Close()
}


// ParseRepository parses the ".db" file passed as reader
func ParseRepository(reader io.Reader, parser Parser) ([]Package, error) {

	// TODO: do actual memory measurings instead of just giving a number out of nowhere
	packages := make([]Package, 0, 10000)

	err := util.ForEachFileInTarGz(reader, func(reader io.Reader) error {
		pkg, err := parser(reader)
		if err != nil {
			return err
		}
		packages = append(packages, pkg)
		return nil
	})

	return packages, err
}
