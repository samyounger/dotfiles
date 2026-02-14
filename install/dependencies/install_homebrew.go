package dependencies

import (
	"fmt"
	"os"
	"os/exec"
)

// HOMEBREW_INSTALL_URL is scoped to this file
const HOMEBREW_INSTALL_URL = "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh"

// InstallHomebrewIfNeeded installs Homebrew if it is not already installed
func InstallHomebrewIfNeeded() {
	if _, err := exec.LookPath("brew"); err != nil {
		fmt.Println("Homebrew not found. Installing Homebrew...")
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("$(curl -fsSL %s)", HOMEBREW_INSTALL_URL))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing Homebrew: %v\n", err)
			os.Exit(1)
		}
	}
}

