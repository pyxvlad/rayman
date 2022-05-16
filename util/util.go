package util

import (
	"errors"
	"os"
	"os/exec"
)

// ConsoleCommand executes a command in a specific workDir
func ConsoleCommand(workDir string, command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir
	return cmd.Run()
}

type ConsoleCommandFunc func(workDir string, name string, arg ...string) error

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

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		} else {
			panic(err)
		}
	}
	return stat.IsDir()
}

func InstallAurPackage(pkg string, console ConsoleCommandFunc) error {
	//	current, err := user.Current()
	//	if err != nil {
	//		return err
	//	}
	//	if current.Name == "root" {
	//		return errors.New("cannot build AUR packages as ROOT")
	//	}
	appCache, err := MakeCacheDir()
	if err != nil {
		return err
	}

	cache := appCache + "/" + pkg

	if DirExists(cache) {
		if err := console(cache, "git", "pull"); err != nil {
			return err
		}
	} else {
		if err := console(appCache, "git", "clone", "https://aur.archlinux.org/"+pkg+".git"); err != nil {
			return err
		}
	}

	return console(cache, "makepkg", "--noconfirm", "-si")
}
