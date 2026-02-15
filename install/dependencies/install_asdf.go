package dependencies

import (
	"fmt"
	"os"
)

var osExit = os.Exit

// checkAsdfInstalled checks if asdf is already installed
func checkAsdfInstalled() {
	if _, err := execLookPath("asdf"); err == nil {
		cmd := execCommand("asdf", "--version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Print("asdf is already installed: ")
		cmd.Run()
		osExit(0)
	}
}

// installAsdf installs asdf using Homebrew
func installAsdf() {
	fmt.Println("Installing asdf using Homebrew...")
	cmd := execCommand("brew", "install", "asdf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing asdf: %v\n", err)
		osExit(1)
	}
	cmd = execCommand("asdf", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Print("asdf installation complete. Version: ")
	cmd.Run()
}

func InstallAsdf() {
	checkAsdfInstalled()
	InstallHomebrew()
	installAsdf()
}
