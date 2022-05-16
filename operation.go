package main

import (
	"fmt"
	"strings"

	"gitlab.com/rayone121/rayman/aurweb"
	"gitlab.com/rayone121/rayman/pacman"
	"gitlab.com/rayone121/rayman/util"
)

type Operation interface {
	Execute() ([]pacman.Package, error)
}

type InstallOperation struct {
	packages         []string
	withConfirmation bool
}

func NewInstallOperation(packages []string, withConfirmation bool) InstallOperation {
	return InstallOperation{packages: packages, withConfirmation: withConfirmation}
}

func (op *InstallOperation) Execute() ([]pacman.Package, error) {

	pac, err := pacman.New()
	if err != nil {
		return nil, err
	}

	pacmanPackages, err := pac.GetAvailablePackages()
	if err != nil {
		return nil, err
	}

	viaPacman := make([]string, 0)
	viaAUR := make([]string, 0)

	for _, i := range op.packages {
		pkgViaPacman := false

		for _, pkg := range pacmanPackages {
			if pkg.Name == i {
				viaPacman = append(viaPacman, i)
				pkgViaPacman = true
				break
			}
		}

		if !pkgViaPacman {
			viaAUR = append(viaAUR, i)
		}
	}

	if len(viaPacman) > 0 {
		err = pacman.Install(viaPacman, op.withConfirmation)
		if err != nil {
			return nil, err
		}
	}

	if len(viaAUR) > 0 {
		results, err := aurweb.Info(viaAUR)
		if err != nil {
			return nil, err
		}

		for _, r := range results {
			util.InstallAurPackage(r.Name)
		}
	}

	return nil, nil
}

type RemoveOperation struct {
	packages         []string
	withConfirmation bool
}

func NewRemoveOperation(packages []string, withConfirmation bool) RemoveOperation {
	return RemoveOperation{packages: packages, withConfirmation: withConfirmation}
}

func (op *RemoveOperation) Execute() ([]pacman.Package, error) {
	packages, err := pacman.GetInstalledPackages()
	if err != nil {
		return nil, err
	}

	for _, i := range op.packages {
		found := false
		for _, pkg := range packages {
			if pkg.Name == i {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("package not found: %s", i)
		}

	}

	return nil, pacman.Remove(op.packages, op.withConfirmation)
}

type SearchField int

const (
	ByName SearchField = iota
	ByNameDesc
	ByMaintainer
	ByDepends
	ByMakeDepends
	ByOptDepends
	ByCheckDepends
)

func (f SearchField) String() string {
	switch f {
	case ByName:
		return "name"
	case ByNameDesc:
		return "name-desc"
	case ByMaintainer:
		return "maintainer"
	case ByDepends:
		return "depends"
	case ByMakeDepends:
		return "makedepends"
	case ByOptDepends:
		return "optdepends"
	case ByCheckDepends:
		return "checkdepends"
	}
	return "unknown"
}

type SearchOperation struct {
	keyword string
	field   string
}

func NewSearchOperation(keyword string, field SearchField) SearchOperation {
	return SearchOperation{keyword: keyword, field: field.String()}
}

func (op *SearchOperation) Execute() ([]pacman.Package, error) {
	pac, err := pacman.New()
	if err != nil {
		return nil, err
	}

	packages := make([]pacman.Package, 0)

	pacmanPackages, err := pac.GetAvailablePackages()
	if err != nil {
		return nil, err
	}

	for _, p := range pacmanPackages {
		if strings.Contains(p.Name, op.keyword) {
			packages = append(packages, p)
		}
	}

	results, err := aurweb.Search(op.field, op.keyword)

	if err != nil {
		return nil, err
	}

	for _, r := range results {
		packages = append(packages, pacman.Package{Name: r.Name, Version: r.Version, Description: r.Description, Repository: "aur"})
	}

	return packages, nil
}
