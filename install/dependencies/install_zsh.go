package dependencies

import (
	"fmt"
	"os"
	"os/exec"
)

// InstallZsh installs zsh using Homebrew if not already installed
func InstallZsh() {
	if _, err := exec.LookPath("zsh"); err == nil {
		fmt.Print("zsh is already installed: ")
		printZshVersion()
		return
	}
	fmt.Println("Installing zsh via Homebrew...")
	cmd := exec.Command("brew", "install", "zsh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing zsh: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("zsh installation complete.")
}

// printZshVersion prints the installed zsh version
func printZshVersion() {
	cmd := exec.Command("zsh", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
