package install

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

var userHomeDir = os.UserHomeDir
var lstat = os.Lstat
var createSymlink = os.Symlink
var runCommand = func(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InstallDotfiles(output io.Writer) {
	if err := applyOSXDefaults(); err != nil {
		fmt.Fprintf(output, "Failed to apply OSX defaults: %v\n", err)
		return
	}

	if err := symlinkHomeFiles(output); err != nil {
		fmt.Fprintf(output, "Failed to create symlinks: %v\n", err)
		return
	}

	fmt.Fprintln(output, "Successful installation.")
}

func applyOSXDefaults() error {
	return runCommand("bash", "-c", "source ./osxdefaults.sh")
}

func symlinkHomeFiles(output io.Writer) error {
	home, err := userHomeDir()
	if err != nil {
		return err
	}

	links := []struct {
		src string
		dst string
	}{
		{filepath.Join(home, ".dotfiles/runcom/.zshrc"), filepath.Join(home, ".zshrc")},
		{filepath.Join(home, ".dotfiles/runcom/.zshenv"), filepath.Join(home, ".zshenv")},
		{filepath.Join(home, ".dotfiles/runcom/.inputrc"), filepath.Join(home, ".inputrc")},
		{filepath.Join(home, ".dotfiles/git/.gitconfig"), filepath.Join(home, ".gitconfig")},
		{filepath.Join(home, ".dotfiles/git/.gitignore_global"), filepath.Join(home, ".gitignore_global")},
	}

	for _, link := range links {
		if _, err := lstat(link.dst); err == nil {
			fmt.Fprintf(output, "Skipping existing file: %s\n", link.dst)
			continue
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if err := createSymlink(link.src, link.dst); err != nil {
			return err
		}
		fmt.Fprintf(output, "Linked %s -> %s\n", link.dst, link.src)
	}

	return nil
}
