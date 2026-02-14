package dependencies

import (
	"fmt"
	"os"
	"os/exec"
)

// InstallGithubCLI installs the GitHub CLI via asdf plugin
func InstallGithubCLI() {
	if isGithubCLIInstalled() {
		fmt.Println("GitHub CLI is already installed.")
		printGithubCLIVersion()
		return
	}
	ensureAsdfInstalled()
	ensureGithubPluginAdded()
	installLatestGithubCLI()
	setGlobalGithubCLIVersion()
	printGithubCLIVersion()
}

// isGithubCLIInstalled checks if the gh CLI is already installed
func isGithubCLIInstalled() bool {
	if _, err := exec.LookPath("gh"); err == nil {
		return true
	}
	return false
}

// ensureAsdfInstalled checks if asdf is installed
func ensureAsdfInstalled() {
	if _, err := exec.LookPath("asdf"); err != nil {
		fmt.Fprintln(os.Stderr, "asdf is not installed. Please install asdf first.")
		os.Exit(1)
	}
}

// ensureGithubPluginAdded checks if the github-cli plugin is added, and adds it if missing
func ensureGithubPluginAdded() {
	pluginListCmd := exec.Command("asdf", "plugin-list")
	pluginListOut, err := pluginListCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing asdf plugins: %v\n", err)
		os.Exit(1)
	}
	if !containsPlugin(string(pluginListOut), "github-cli") {
		fmt.Println("Adding github-cli plugin to asdf...")
		if err := addGithubPlugin(); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding github-cli plugin: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("github-cli plugin already added to asdf.")
	}
}

// addGithubPlugin tries to add the plugin using SSH, then HTTPS with GIT_ASKPASS fallback
func addGithubPlugin() error {
	var addCmd = exec.Command("asdf", "plugin-add", "github-cli")
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr
	addCmd.Env = append(os.Environ(), "GIT_ASKPASS=/bin/echo")
	return addCmd.Run()
}

// installLatestGithubCLI installs the latest github-cli version via asdf
func installLatestGithubCLI() {
	fmt.Println("Installing latest github-cli via asdf...")
	installCmd := exec.Command("asdf", "install", "github-cli", "latest")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing github-cli: %v\n", err)
		os.Exit(1)
	}
}

// setGlobalGithubCLIVersion sets the global github-cli version to latest
func setGlobalGithubCLIVersion() {
	fmt.Println("Setting github-cli global version to latest...")
	globalCmd := exec.Command("asdf", "global", "github-cli", "latest")
	globalCmd.Stdout = os.Stdout
	globalCmd.Stderr = os.Stderr
	if err := globalCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting github-cli global version: %v\n", err)
		os.Exit(1)
	}
}

// printGithubCLIVersion prints the installed github-cli version
func printGithubCLIVersion() {
	versionCmd := exec.Command("gh", "--version")
	versionCmd.Stdout = os.Stdout
	versionCmd.Stderr = os.Stderr
	fmt.Print("GitHub CLI installation complete. Version: ")
	versionCmd.Run()
}

// containsPlugin checks if the plugin is in the asdf plugin list output
func containsPlugin(list, plugin string) bool {
	for _, line := range splitLines(list) {
		if line == plugin {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
