package pacman_test

import (
	"testing"

	"gitlab.com/rayone121/rayman/pacman"
)

func TestGetInstalledPackages(t *testing.T) {
	
	packages ,err := pacman.GetInstalledPackages()
	if err != nil {
		t.Fatalf("error: %e", err)
	}

	found := false
	for _, p := range packages {
		if p.Name == "linux" {
			found = true
		}
	}

	if !found {
		t.Fatal("something went wrong... couldn't find 'linux' package...")
	}


}

// NOTE: it assumes an Arch Linux system with the linux package installed
func TestIsPackageInstalled(t *testing.T) {

	installed, err := pacman.IsPackageInstalled("linux")
	if err != nil {
		t.Fatalf("error: %e", err);
	}

	if !installed {
		t.Fatal("seems like linux isn't installed...")
	}
}

