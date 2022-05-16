package util

import (
	"errors"
	"os"
	"os/exec"
)

// NewConsoleCommand uses exec.Command to create a new *exec.Cmd and also sets it to use standard I/O
func NewConsoleCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd

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
		cmd := NewConsoleCommand("git", "pull")
		cmd.Dir = cache + "/" + pkg
		if err != nil {
			return err
		}
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else {
		err := os.RemoveAll(cache + "/" + pkg)
		if err != nil {
			return err
		}

		cmd := NewConsoleCommand("git", "clone", "https://aur.archlinux.org/"+pkg+".git")
		cmd.Dir = cache
		if err != nil {
			return err
		}

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	cmd := NewConsoleCommand("makepkg", "-si")
	cmd.Dir = cache + "/" + pkg
	if err != nil {
		return err
	}

	return cmd.Run()
}
