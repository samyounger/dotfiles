package dependencies

import (
	"fmt"
	"os"
	"os/exec"
)

// HOMEBREW_INSTALL_URL is scoped to this file
const HOMEBREW_INSTALL_URL = "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh"

// isHomebrewInstalled checks if Homebrew is already installed
func isHomebrewInstalled() bool {
	if _, err := execLookPath("brew"); err == nil {
		return true
	}
	return false
}

// printHomebrewVersion prints the installed Homebrew version
func printHomebrewVersion() {
	cmd := exec.Command("brew", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// InstallHomebrew installs Homebrew if it is not already installed
func InstallHomebrew() {
	if isHomebrewInstalled() {
		fmt.Print("Homebrew is already installed:")
		printHomebrewVersion()
		return
	}
	fmt.Println("Homebrew not found. Installing Homebrew...")
	cmd := execCommand("/bin/bash", "-c", fmt.Sprintf("$(curl -fsSL %s)", HOMEBREW_INSTALL_URL))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing Homebrew: %v\n", err)
		os.Exit(1)
	}
}
