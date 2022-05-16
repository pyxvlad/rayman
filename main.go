package main

import (
	"fmt"
	"gitlab.com/rayone121/rayman/pacman"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("insufficient arguments")
		return
	}

	pkg := os.Args[2]
	switch os.Args[1] {
	case "install":
		{
			op := pacman.NewInstallOperation(os.Args[2:], true)
			_, err := op.Execute()
			if err != nil {
				panic(err)
			}

		}
	case "remove":
		{
			op := pacman.NewRemoveOperation(os.Args[2:], true)
			_, err := op.Execute()
			if err != nil {
				panic(err)
			}

		}

	case "search":
		{
			op := pacman.NewSearchOperation(pkg, pacman.ByName)
			results, err := op.Execute()
			if err != nil {
				panic(err)
			}

			for _, pkg := range results {
				fmt.Printf("\n%s/%s %s", pkg.Repository, pkg.Name, pkg.Version)
				installed, err := pacman.IsPackageInstalled(pkg.Name)
				if err != nil {
					panic(err)
				}
				if installed {
					fmt.Print("[installed]")
				}

				fmt.Print("\n\t", pkg.Description)

			}
			fmt.Print('\n')

		}
	}
}
