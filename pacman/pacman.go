package pacman

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/rayone121/rayman/util"
)

type Pacman struct {
	dbPath string
	repos []string
}

// New creates a new Pacman instance with the specific dataDir
// Most probably you want to set dataDir to "/var/lib/pacman/sync"
func New(dataDir string) (Pacman, error) {

	var pacman Pacman
	pacman.dbPath = dataDir
	entries, err := os.ReadDir(pacman.dbPath)
	if err != nil {
		return pacman, err
	}

	pacman.repos = make([]string, 0)

	for _, e := range entries {
		if e.Type().IsRegular() && strings.HasSuffix(e.Name(), ".db") {
			pacman.repos = append(pacman.repos, strings.TrimSuffix(e.Name(), ".db"))
		}

	}

	return pacman, nil
}

func (p *Pacman) GetAvailablePackages() ([]Package, error) {
	packages := make([]Package, 0, 100000)
	for _, r := range p.repos {
		repoPackages, err := ParseRepositoryFile(p.dbPath + "/" + r + ".db", FromDescReader)
		if err != nil {
			return nil, err
		}

		for _, i := range repoPackages {
			packages = append(packages, i)
		}

	}
	return packages, nil
}

func (p *Pacman) GetRepositoryPackages(repo string) ([]Package, error) {
	found := false
	for _, r := range p.repos {
		if r == repo {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("repository %s doesn't exist", repo)
	}

	repoPackages, err := ParseRepositoryFile(p.dbPath + "/" + repo + ".db", FromDescReader)
	if err != nil {
		return nil, err
	}
	return repoPackages, err
}

type Package struct {
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Repository  string `json:"repository,omitempty"`
}

func GetInstalledPackages() ([]Package, error) {
	output, err := exec.Command("pacman", "-Q").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.SplitN(string(output), "\n", -1)

	packages := make([]Package, 1000)
	for _, l := range lines {
		split := strings.Split(l, " ")
		// TODO: replace this ugly TEMPORARY hack, added it because it was index-out-of-bounding lol
		if len(split) == 2 {
			packages = append(packages, Package{Name: split[0], Version: split[1]})
		}
	}

	return packages, nil
}

func Install(pkgs []string, confirm bool) error {
	args := []string{"pacman", "--sync"}
	if !confirm {
		args = append(args, "--noconfirm")
	}

	for _, pkg := range pkgs {
		args = append(args, pkg)
	}
	return util.ConsoleCommand(".", "sudo", args...)
}

func Remove(pkgs []string, confirm bool) error {
	args := []string{"pacman", "--remove"}
	if !confirm {
		args = append(args, "--noconfirm")
	}

	for _, pkg := range pkgs {
		args = append(args, pkg)
	}
	return util.ConsoleCommand(".", "sudo", args...)
}

func IsPackageInstalled(pkg string) (bool, error) {
	packages, err := GetInstalledPackages()
	if err != nil {
		return false, err
	}
	for _, p := range packages {
		if p.Name == pkg {
			return true, nil
		}
	}
	return false, nil
}
