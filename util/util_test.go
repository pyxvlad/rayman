package util_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"gitlab.com/rayone121/rayman/util"
)

func TestConsoleCommand(t *testing.T) {
	tmp := t.TempDir()
	if err := util.ConsoleCommand(tmp, "ls"); err != nil {
		t.Fatal(err)
	}
}

func TestMakeCacheDir_NoXDGNoHome(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "")
	defer func() { recover() }()
	util.MakeCacheDir()
	t.Fatal("didn't panic")
}

func TestMakeCacheDir_NoXDG(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", tmp)
	cache, err := util.MakeCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if cache != (tmp + "/.cache/rayman") {
		t.Fatalf("cache should be %s/.config/rayman instead it's %s", tmp, cache)
	}
}

func TestMakeCacheDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)
	cache, err := util.MakeCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if cache != tmp+"/rayman" {
		t.Fatalf("cache should be $XDG_CACHE_HOME/rayman instead it's %s", cache)
	}
}

func TestDirExists_InvalidDir(t *testing.T) {
	tmp := t.TempDir()
	exists := util.DirExists(tmp + "/invalid_dir")
	if exists {
		t.Fatal("invalid directory yet it looks like it exists")
	}
}

func TestDirExists_NotADirectory(t *testing.T) {
	tmp := t.TempDir()
	filename := tmp + "/not_a_directory"
	os.Create(filename)
	defer func() { recover() }()
	util.DirExists(filename + "/invalid_dir")
	t.Fatal("expected panic")
}

func TestDirExists(t *testing.T) {
	tmp := t.TempDir()
	if !util.DirExists(tmp) {
		t.Fatal("directory should exist, yet it doesn't")
	}
}

type FakeConsole struct {
	t       *testing.T
	err     error
	workDir string
	command string
	args    []string
}

func (f FakeConsole) Console(workdir string, command string, args ...string) error {
	if f.workDir != workdir {
		f.t.Fatalf("expected workdir %s got %s", f.workDir, workdir)
	}
	if f.command != command {
		f.t.Fatalf("expected command %s got %s", f.command, command)
	}
	if !reflect.DeepEqual(f.args, args) {
		f.t.Fatalf("expected args %#v got %#v", f.args, args)
	}
	return f.err
}

func NewFakeConsole(t *testing.T, err error, workdir string, command string, args ...string) FakeConsole {
	return FakeConsole{t, err, workdir, command, args}
}

type FakeConsoleSeries struct {
	fakes chan FakeConsole
	t     *testing.T
}

func (f FakeConsoleSeries) Console(w string, c string, a ...string) error {
	select {
	case fc := <-f.fakes:
		return fc.Console(w, c, a...)
	default:
		f.t.Fatal("no more console fakes")
	}
	panic("unreachable")
}

func (f FakeConsoleSeries) Add(err error, w string, c string, a ...string) {
	f.fakes <- NewFakeConsole(f.t, err, w, c, a...)
}

func NewFakeConsoleSeries(t *testing.T, num int) FakeConsoleSeries {
	return FakeConsoleSeries{t: t, fakes: make(chan FakeConsole, num)}
}

func TestInstallAurPackage_InvalidCacheDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)
	if _, err := os.Create(tmp + "/rayman"); err != nil {
		t.Fatal(err)
	}

	if err := util.InstallAurPackage("rayman", func(workDir, name string, arg ...string) error { return nil }); err == nil {
		t.Fatal("expected failure didn't happen")
	}
}

func TestInstallAurPackage_CacheNotADirectory(t *testing.T) {
	tmp := t.TempDir()
	filename := tmp + "/not_a_dir"
	t.Setenv("XDG_CACHE_HOME", filename+"/f")

	if _, err := os.Create(filename); err != nil {
		t.Fatal(err)
	}

	if err := util.InstallAurPackage(
		"rayman", func(workDir, name string, arg ...string) error {
			return nil
		}); err == nil {
		t.Fatal("expected failure didn't happen")
	}
}


func TestInstallAurPackage_FailGitClone(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)

	fakes := NewFakeConsoleSeries(t, 1)
	gitErr := errors.New("git fail")
	fakes.Add(gitErr, tmp+"/rayman", "git", "clone", "https://aur.archlinux.org/rayman.git")
	if err := util.InstallAurPackage("rayman", fakes.Console); !errors.Is(err, gitErr) {
		t.Fatalf("expected error %e got %e", gitErr, err)
	}
}

func TestInstallAurPackage_FailCleanDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)

	os.MkdirAll(tmp + "/rayman", 0777)
	f, err := os.Create(tmp + "/rayman/rayman")
	if err != nil {
		t.Fatal(err)
	} 
	if err := f.Chmod(0); err != nil {
		t.Fatal(err)
	}

	if err := util.InstallAurPackage("rayman", func(workDir, name string, arg ...string) error {return nil}); err != nil {
		t.Fatalf("got %e", err)
	}
}

func TestInstallAurPackage_FailGitPull(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)

	if err := os.MkdirAll(tmp + "/rayman/rayman", 0777); err != nil {
		t.Fatalf("internal error: %e", err)
	}

	fakes := NewFakeConsoleSeries(t, 1)
	gitErr := errors.New("git fail")
	fakes.Add(gitErr, tmp+"/rayman/rayman", "git", "pull")
	if err := util.InstallAurPackage("rayman", fakes.Console); !errors.Is(err, gitErr) {
		t.Fatalf("expected error %e got %e", gitErr, err)
	}
}

func TestInstallAurPackage_Fresh(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)

	fakes := NewFakeConsoleSeries(t, 2)
	fakes.Add(nil, tmp+"/rayman", "git", "clone", "https://aur.archlinux.org/rayman.git")
	fakes.Add(nil, tmp+"/rayman/rayman", "makepkg", "--noconfirm", "-si")
	util.InstallAurPackage("rayman", fakes.Console)
}
