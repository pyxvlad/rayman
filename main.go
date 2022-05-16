package main

import (
	"errors"
	"fmt"
	"gitlab.com/rayone121/rayman/aurweb"
	"os"
	"os/exec"
	"strings"
)

func IsPackageInstalled(pkg string) (bool, error) {
	output, err := exec.Command("pacman", "-Q").Output()
	if err != nil {
		return false, err
	}

	lines := strings.SplitN(string(output), "\n", -1)

	for _, l := range lines {
		split := strings.Split(l, " ")
		if split[0] == pkg {
			return true, nil
		}
	}
	return false, nil
}

func MakeCacheDir() (string, error) {
	xdgCache := os.Getenv("XDG_CACHE_HOME")
	if xdgCache == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		xdgCache = home + "/.cache"
	}

	cache := xdgCache + "/rayman"

	return cache, os.MkdirAll(cache, 0777)
}
func DirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}
	if stat.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func InstallAurPackage(pkg string) error {
	cache, err := MakeCacheDir()
	if err != nil {
		return err
	}

	exists, err := DirExists(cache + "/" + pkg)
	if err != nil {
		return err
	}

	if exists {
		cmd := exec.Command("git", "pull")
		cmd.Dir = cache + "/" + pkg
		if err != nil {
			return err
		}
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stdin = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		err := os.RemoveAll(cache + "/" + pkg)
		if err != nil {
			return err
		}

		cmd := exec.Command("git", "clone", "https://aur.archlinux.org/"+pkg+".git")
		cmd.Dir = cache
		if err != nil {
			return err
		}
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stdin = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("makepkg", "-si")
	cmd.Dir = cache + "/" + pkg
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("insufficient arguments")
		return
	}

	pkg := os.Args[2]
	switch os.Args[1] {
	case "install":
		{
			installed, err := IsPackageInstalled(pkg)
			if err != nil {
				return
			}
			if installed {
				fmt.Printf("It looks like package %s is already installed... Do you want to continue? [Y/n]", pkg)
				var answer rune
				_, err := fmt.Scan(&answer)
				if err != nil {
					panic(err)
				}
				if answer == 'n' || answer == 'N' {
					fmt.Println("Ok... aborting...")
					os.Exit(1)
				}
			}

			cmd := exec.Command("sudo", "pacman", "--sync", "--noconfirm", pkg)

			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				fmt.Println("Pacman failed to install package, trying to find the package in the AUR")
				info, err := aurweb.Info([]string{pkg})
				if err != nil {
					panic(err)
				}
				if len(info) == 0 {
					fmt.Println("No AUR package found... Aborting...")
					os.Exit(1)
				}

				err = InstallAurPackage(pkg)
				if err != nil {
					fmt.Println(err)
				}
			}

		}
	case "remove":
		{
			installed, err := IsPackageInstalled(pkg)
			if err != nil {
				panic(err)
			}
			if !installed {
				fmt.Println("Package is not installed!")
			}

			cmd := exec.Command("sudo", "pacman", "-R", pkg)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				panic(err)
			}

		}

	case "search":
		{
			cmd := exec.Command("pacman", "-Ss", pkg)
			cmd.Stdout = os.Stdout

			err := cmd.Run()
			if err != nil {
				fmt.Printf("pacman: %s", err)
			}

			results, err := aurweb.Search("name", pkg)
			if err != nil {
				panic(err)
			}

			for _, r := range results {
				fmt.Printf("\naur/%s %s", r.Name, r.Version)
				installed, err := IsPackageInstalled(r.Name)
				if err != nil {
					panic(err)
				}
				if installed {
					fmt.Print("[installed]")
				}

				fmt.Print("\n\t", r.Description)

			}
			fmt.Print('\n')

		}
	}
}
