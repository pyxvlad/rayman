package pacman

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/rayone121/rayman/aurweb"
	"gitlab.com/rayone121/rayman/util"
)

type Operation interface {
	Execute() ([]Package, error)
}

type operationJsonHelper struct {
	Type     string   `json:"type,omitempty"`
	Packages []string `json:"packages,omitempty"`
}

type InstallOperation struct {
	packages         []string
	withConfirmation bool
}

func NewInstallOperation(packages []string, withConfirmation bool) InstallOperation {
	return InstallOperation{packages: packages, withConfirmation: withConfirmation}
}

func (op InstallOperation) Execute() ([]Package, error) {

	pac, err := New()
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
		err = Install(viaPacman, op.withConfirmation)
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
			err := util.InstallAurPackage(r.Name)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (op InstallOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(operationJsonHelper{Type: "install", Packages: op.packages})
}

type RemoveOperation struct {
	packages         []string
	withConfirmation bool
}

func NewRemoveOperation(packages []string, withConfirmation bool) RemoveOperation {
	return RemoveOperation{packages: packages, withConfirmation: withConfirmation}
}

func (op RemoveOperation) Execute() ([]Package, error) {
	packages, err := GetInstalledPackages()
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

	return nil, Remove(op.packages, op.withConfirmation)
}

func (op RemoveOperation) MarshalJSON() ([]byte, error) {

	return json.Marshal(operationJsonHelper{Type: "remove", Packages: op.packages})
}

type SearchField string

const (
	ByName         SearchField = "name"
	ByNameDesc                 = "name-desc"
	ByMaintainer               = "maintainer"
	ByDepends                  = "depends"
	ByMakeDepends              = "makedepends"
	ByOptDepends               = "optdepends"
	ByCheckDepends             = "checkdepends"
)

type SearchOperation struct {
	keyword string
	field   SearchField
}

func NewSearchOperation(keyword string, field SearchField) SearchOperation {
	return SearchOperation{keyword: keyword, field: field}
}

func (op SearchOperation) Execute() ([]Package, error) {
	pac, err := New()
	if err != nil {
		return nil, err
	}

	packages := make([]Package, 0)

	pacmanPackages, err := pac.GetAvailablePackages()
	if err != nil {
		return nil, err
	}

	for _, p := range pacmanPackages {
		if strings.Contains(p.Name, op.keyword) {
			packages = append(packages, p)
		}
	}

	results, err := aurweb.Search(string(op.field), op.keyword)

	if err != nil {
		return nil, err
	}

	for _, r := range results {
		packages = append(packages, Package{Name: r.Name, Version: r.Version, Description: r.Description, Repository: "aur"})
	}

	return packages, nil
}
