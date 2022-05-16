package pacman

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

const descFileFields = 16

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

func ParseRepositoryFile(repoName string) ([]Package, error) {
	file, err := os.Open("/var/lib/pacman/sync/" + repoName + ".db")
	if err != nil {
		return nil, err
	}

	packages, err := ParseRepository(file, repoName)
	if err != nil {
		if err := file.Close(); err != nil {
			panic(err)
		}
		return nil, err
	}

	return packages, file.Close()
}

func ParseRepository(reader io.Reader, repoName string) ([]Package, error) {

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(gzReader)

	// TODO: do actual memory measurings instead of just giving a number out of nowhere
	packages := make([]Package, 0, 10000)

	header, err := tarReader.Next()
	for err != io.EOF {
		if header.Typeflag == tar.TypeReg {
			pkg, err := FromDescReader(tarReader)
			if err != nil {
				return nil, err
			}

			pkg.Repository = repoName
			packages = append(packages, pkg)
		}
		header, err = tarReader.Next()
	}
	return packages, nil
}
