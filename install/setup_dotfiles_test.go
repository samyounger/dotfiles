package install

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInstallDotfiles_Success(t *testing.T) {
	defer resetDotfilesMocks()

	userHomeDir = func() (string, error) { return "/tmp/home", nil }
	lstat = func(name string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	created := map[string]string{}
	copied := map[string]string{}
	createSymlink = func(oldname, newname string) error {
		created[newname] = oldname
		return nil
	}
	readFile = func(name string) ([]byte, error) {
		if name == "/tmp/home/.dotfiles/runcom/.tool-versions" {
			return []byte("golang 1.26.1\n"), nil
		}
		return nil, errors.New("unexpected read")
	}
	writeFile = func(name string, data []byte, perm os.FileMode) error {
		copied[name] = string(data)
		return nil
	}

	out := &bytes.Buffer{}
	InstallDotfiles(out)

	expected := map[string]string{
		"/tmp/home/.zshrc":            "/tmp/home/.dotfiles/runcom/.zshrc",
		"/tmp/home/.zshenv":           "/tmp/home/.dotfiles/runcom/.zshenv",
		"/tmp/home/.inputrc":          "/tmp/home/.dotfiles/runcom/.inputrc",
		"/tmp/home/.gitconfig":        "/tmp/home/.dotfiles/git/.gitconfig",
		"/tmp/home/.gitignore_global": "/tmp/home/.dotfiles/git/.gitignore_global",
	}

	if len(created) != len(expected) {
		t.Fatalf("expected %d symlinks, got %d", len(expected), len(created))
	}
	for dst, src := range expected {
		if created[dst] != src {
			t.Fatalf("expected link %s -> %s, got %s", dst, src, created[dst])
		}
	}
	if copied["/tmp/home/.tool-versions"] != "golang 1.26.1\n" {
		t.Fatalf("expected .tool-versions to be copied, got %q", copied["/tmp/home/.tool-versions"])
	}
	if !strings.Contains(out.String(), "Successful installation.") {
		t.Fatalf("expected success message, got: %s", out.String())
	}
}

func TestInstallDotfiles_HomeDirFailure(t *testing.T) {
	defer resetDotfilesMocks()

	userHomeDir = func() (string, error) { return "", errors.New("boom") }

	out := &bytes.Buffer{}
	InstallDotfiles(out)

	if !strings.Contains(out.String(), "Failed to create symlinks") {
		t.Fatalf("expected symlink error message, got: %s", out.String())
	}
}

func TestInstallDotfiles_SymlinkFailure(t *testing.T) {
	defer resetDotfilesMocks()

	userHomeDir = func() (string, error) { return "/tmp/home", nil }
	lstat = func(name string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	createSymlink = func(oldname, newname string) error {
		if filepath.Base(newname) == ".zshrc" {
			return errors.New("cannot link")
		}
		return nil
	}

	out := &bytes.Buffer{}
	InstallDotfiles(out)

	if !strings.Contains(out.String(), "Failed to create symlinks") {
		t.Fatalf("expected symlink error message, got: %s", out.String())
	}
}

func TestSymlinkHomeFiles_SkipsExisting(t *testing.T) {
	defer resetDotfilesMocks()

	userHomeDir = func() (string, error) { return "/tmp/home", nil }
	lstat = func(name string) (fs.FileInfo, error) {
		if strings.HasSuffix(name, ".zshrc") {
			return fakeFileInfo{}, nil
		}
		return nil, fs.ErrNotExist
	}
	created := 0
	wroteToolVersions := false
	createSymlink = func(oldname, newname string) error {
		created++
		return nil
	}
	readFile = func(name string) ([]byte, error) { return []byte("golang 1.26.1\n"), nil }
	writeFile = func(name string, data []byte, perm os.FileMode) error {
		if name == "/tmp/home/.tool-versions" {
			wroteToolVersions = true
		}
		return nil
	}

	out := &bytes.Buffer{}
	if err := symlinkHomeFiles(out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Skipping existing file: /tmp/home/.zshrc") {
		t.Fatalf("expected skip message, got: %s", out.String())
	}
	if created != 4 {
		t.Fatalf("expected 4 symlinks to be created, got %d", created)
	}
	if !wroteToolVersions {
		t.Fatalf("expected .tool-versions to be written")
	}
}

func TestSymlinkHomeFiles_CopyFailure(t *testing.T) {
	defer resetDotfilesMocks()

	userHomeDir = func() (string, error) { return "/tmp/home", nil }
	lstat = func(name string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	createSymlink = func(oldname, newname string) error { return nil }
	readFile = func(name string) ([]byte, error) {
		return nil, errors.New("cannot read")
	}

	out := &bytes.Buffer{}
	err := symlinkHomeFiles(out)
	if err == nil || !strings.Contains(err.Error(), "cannot read") {
		t.Fatalf("expected copy failure, got: %v", err)
	}
}

func resetDotfilesMocks() {
	userHomeDir = os.UserHomeDir
	lstat = os.Lstat
	createSymlink = os.Symlink
	readFile = os.ReadFile
	writeFile = os.WriteFile
	runCommand = func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

type fakeFileInfo struct{}

func (fakeFileInfo) Name() string       { return "" }
func (fakeFileInfo) Size() int64        { return 0 }
func (fakeFileInfo) Mode() fs.FileMode  { return 0 }
func (fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (fakeFileInfo) IsDir() bool        { return false }
func (fakeFileInfo) Sys() any           { return nil }
