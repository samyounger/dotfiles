package dependencies

import (
	"fmt"
	"os"
	"os/exec"
)

// checkGitInstalled checks if git is already installed
func checkGitInstalled() {
	if _, err := exec.LookPath("git"); err == nil {
		cmd := exec.Command("git", "--version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Print("git is already installed: ")
		cmd.Run()
		os.Exit(0)
	}
}

// installGit installs git using Homebrew
func installGit() {
	fmt.Println("Installing git using Homebrew...")
	cmd := exec.Command("brew", "install", "git")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing git: %v\n", err)
		os.Exit(1)
	}
	cmd = exec.Command("git", "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Print("git installation complete. Version: ")
	cmd.Run()
}

func InstallGit() {
	checkGitInstalled()
	InstallHomebrewIfNeeded()
	installGit()
}
