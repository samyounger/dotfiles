package dependencies

import (
	"fmt"
	"os"
)

// InstallZsh installs zsh using Homebrew if not already installed
func InstallZsh() {
	if _, err := execLookPath("zsh"); err == nil {
		fmt.Print("zsh is already installed: ")
		printZshVersion()
		return
	}
	fmt.Println("Installing zsh via Homebrew...")
	cmd := execCommand("brew", "install", "zsh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing zsh: %v\n", err)
		osExit(1)
	}
	fmt.Println("zsh installation complete.")
}

// printZshVersion prints the installed zsh version
func printZshVersion() {
	cmd := execCommand("zsh", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
