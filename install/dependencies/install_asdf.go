package dependencies

import (
	"fmt"
	"os"
	"os/exec"
)

// checkAsdfInstalled checks if asdf is already installed
func checkAsdfInstalled() {
	if _, err := exec.LookPath("asdf"); err == nil {
		cmd := exec.Command("asdf", "--version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Print("asdf is already installed: ")
		cmd.Run()
		os.Exit(0)
	}
}

// installAsdf installs asdf using Homebrew
func installAsdf() {
	fmt.Println("Installing asdf using Homebrew...")
	cmd := exec.Command("brew", "install", "asdf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing asdf: %v\n", err)
		os.Exit(1)
	}
	cmd = exec.Command("asdf", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Print("asdf installation complete. Version: ")
	cmd.Run()
}

func InstallAsdf() {
	checkAsdfInstalled()
	InstallHomebrewIfNeeded()
	installAsdf()
}
