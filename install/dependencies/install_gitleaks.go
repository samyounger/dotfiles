package dependencies

import (
	"fmt"
	"os"
	"path/filepath"
)

// InstallGitleaks installs gitleaks via asdf.
func InstallGitleaks() {
	if isGitleaksInstalled() {
		fmt.Println("gitleaks is already installed.")
		printGitleaksVersion()
		return
	}
	ensureAsdfInstalled()
	ensureGitleaksPluginAdded()
	version, err := desiredGitleaksVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error determining gitleaks version: %v\n", err)
		osExit(1)
	}
	installGitleaksVersion(version)
	setGlobalGitleaksVersion(version)
	printGitleaksVersion()
}

func isGitleaksInstalled() bool {
	if _, err := execLookPath("gitleaks"); err == nil {
		return true
	}
	return false
}

func ensureGitleaksPluginAdded() {
	pluginListCmd := execCommand("asdf", "plugin-list")
	pluginListOut, err := pluginListCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing asdf plugins: %v\n", err)
		osExit(1)
	}
	if !containsPlugin(string(pluginListOut), "gitleaks") {
		fmt.Println("Adding gitleaks plugin to asdf...")
		addCmd := execCommand("asdf", "plugin-add", "gitleaks")
		addCmd.Stdout = os.Stdout
		addCmd.Stderr = os.Stderr
		if err := addCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding gitleaks plugin: %v\n", err)
			osExit(1)
		}
	} else {
		fmt.Println("gitleaks plugin already added to asdf.")
	}
}

func installGitleaksVersion(version string) {
	fmt.Printf("Installing gitleaks %s via asdf...\n", version)
	installCmd := execCommand("asdf", "install", "gitleaks", version)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing gitleaks: %v\n", err)
		osExit(1)
	}
}

func setGlobalGitleaksVersion(version string) {
	fmt.Printf("Setting gitleaks global version to %s...\n", version)
	globalCmd := execCommand("asdf", "global", "gitleaks", version)
	globalCmd.Stdout = os.Stdout
	globalCmd.Stderr = os.Stderr
	if err := globalCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting gitleaks global version: %v\n", err)
		osExit(1)
	}
}

func printGitleaksVersion() {
	versionCmd := execCommand("gitleaks", "version")
	versionCmd.Stdout = os.Stdout
	versionCmd.Stderr = os.Stderr
	versionCmd.Run()
}

func desiredGitleaksVersion() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	toolVersionsPath := filepath.Join(home, ".dotfiles", "runcom", ".tool-versions")
	content, err := os.ReadFile(toolVersionsPath)
	if err != nil {
		return "", err
	}

	for _, line := range splitLines(string(content)) {
		if len(line) == 0 {
			continue
		}
		fields := splitFields(line)
		if len(fields) >= 2 && fields[0] == "gitleaks" {
			return fields[1], nil
		}
	}

	return "", fmt.Errorf("gitleaks entry not found in %s", toolVersionsPath)
}

func splitFields(s string) []string {
	var fields []string
	fieldStart := -1
	for i, r := range s {
		if r == ' ' || r == '\t' {
			if fieldStart != -1 {
				fields = append(fields, s[fieldStart:i])
				fieldStart = -1
			}
			continue
		}
		if fieldStart == -1 {
			fieldStart = i
		}
	}
	if fieldStart != -1 {
		fields = append(fields, s[fieldStart:])
	}
	return fields
}
