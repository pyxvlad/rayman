package main

import (
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

func main() {
	if len(os.Args) < 3 {
		fmt.Println("insufficient arguments")
		return
	}

	pkg := os.Args[2]
	switch os.Args[1] {
	case "install":
		{

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
			cmd.Stderr = os.Stdin
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
				panic(err)
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

		}
	}
}
