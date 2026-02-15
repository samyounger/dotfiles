package dependencies

import (
	"fmt"
	"os"
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
	if _, err := execLookPath("gh"); err == nil {
		return true
	}
	return false
}

// ensureAsdfInstalled checks if asdf is installed
func ensureAsdfInstalled() {
	if _, err := execLookPath("asdf"); err != nil {
		fmt.Fprintln(os.Stderr, "asdf is not installed. Please install asdf first.")
		osExit(1)
	}
}

// ensureGithubPluginAdded checks if the github-cli plugin is added, and adds it if missing
func ensureGithubPluginAdded() {
	pluginListCmd := execCommand("asdf", "plugin-list")
	pluginListOut, err := pluginListCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing asdf plugins: %v\n", err)
		osExit(1)
	}
	if !containsPlugin(string(pluginListOut), "github-cli") && !containsPlugin(string(pluginListOut), "gh") {
		fmt.Println("Adding github-cli plugin to asdf...")
		if err := addGithubPlugin(); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding github-cli plugin: %v\n", err)
			osExit(1)
		}
	} else {
		fmt.Println("github-cli plugin already added to asdf.")
	}
}

// addGithubPlugin tries to add the plugin using SSH, then HTTPS with GIT_ASKPASS fallback
func addGithubPlugin() error {
	sshCmd := execCommand("asdf", "plugin-add", "github-cli", "git@github.com:kennyp/asdf-gh.git")
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr
	if err := sshCmd.Run(); err == nil {
		return nil
	}

	httpsCmd := execCommand("asdf", "plugin-add", "github-cli", "https://github.com/kennyp/asdf-gh.git")
	httpsCmd.Stdout = os.Stdout
	httpsCmd.Stderr = os.Stderr
	httpsCmd.Env = append(os.Environ(), "GIT_ASKPASS=/bin/echo")
	return httpsCmd.Run()
}

// installLatestGithubCLI installs the latest github-cli version via asdf
func installLatestGithubCLI() {
	fmt.Println("Installing latest github-cli via asdf...")
	installCmd := execCommand("asdf", "install", "github-cli", "latest")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing github-cli: %v\n", err)
		osExit(1)
	}
}

// setGlobalGithubCLIVersion sets the global github-cli version to latest
func setGlobalGithubCLIVersion() {
	fmt.Println("Setting github-cli global version to latest...")
	globalCmd := execCommand("asdf", "global", "github-cli", "latest")
	globalCmd.Stdout = os.Stdout
	globalCmd.Stderr = os.Stderr
	if err := globalCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting github-cli global version: %v\n", err)
		osExit(1)
	}
}

// printGithubCLIVersion prints the installed github-cli version
func printGithubCLIVersion() {
	versionCmd := execCommand("gh", "--version")
	versionCmd.Stdout = os.Stdout
	versionCmd.Stderr = os.Stderr
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
