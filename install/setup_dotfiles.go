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
var readFile = os.ReadFile
var writeFile = os.WriteFile
var runCommand = func(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InstallDotfiles(output io.Writer) {
	if err := symlinkHomeFiles(output); err != nil {
		fmt.Fprintf(output, "Failed to create symlinks: %v\n", err)
		return
	}

	fmt.Fprintln(output, "Successful installation.")
}

func symlinkHomeFiles(output io.Writer) error {
	home, err := userHomeDir()
	if err != nil {
		return err
	}

	type installKind int
	const (
		installSymlink installKind = iota
		installCopy
	)

	links := []struct {
		src  string
		dst  string
		kind installKind
	}{
		{filepath.Join(home, ".dotfiles/runcom/.zshrc"), filepath.Join(home, ".zshrc"), installSymlink},
		{filepath.Join(home, ".dotfiles/runcom/.zshenv"), filepath.Join(home, ".zshenv"), installSymlink},
		{filepath.Join(home, ".dotfiles/runcom/.inputrc"), filepath.Join(home, ".inputrc"), installSymlink},
		{filepath.Join(home, ".dotfiles/git/.gitconfig"), filepath.Join(home, ".gitconfig"), installSymlink},
		{filepath.Join(home, ".dotfiles/git/.gitignore_global"), filepath.Join(home, ".gitignore_global"), installSymlink},
		{filepath.Join(home, ".dotfiles/runcom/.tool-versions"), filepath.Join(home, ".tool-versions"), installCopy},
	}

	for _, link := range links {
		if _, err := lstat(link.dst); err == nil {
			fmt.Fprintf(output, "Skipping existing file: %s\n", link.dst)
			continue
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		switch link.kind {
		case installSymlink:
			if err := createSymlink(link.src, link.dst); err != nil {
				return err
			}
			fmt.Fprintf(output, "Linked %s -> %s\n", link.dst, link.src)
		case installCopy:
			contents, err := readFile(link.src)
			if err != nil {
				return err
			}
			if err := writeFile(link.dst, contents, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(output, "Copied %s -> %s\n", link.dst, link.src)
		}
	}

	return nil
}
