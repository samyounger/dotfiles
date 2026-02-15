package install

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/syounger-ut/dotfiles/install/application_settings"
)

var runCommandOutput = func(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func SetupOSXDefaults(input io.Reader, output io.Writer) {
	fmt.Fprintln(output, "OSX Defaults Options:")
	fmt.Fprintln(output, "1) Setup MacOS settings")
	fmt.Fprintln(output, "0) Exit")
	fmt.Fprint(output, "$ ")

	scanner := bufio.NewScanner(input)
	if !scanner.Scan() {
		fmt.Fprintln(output, "Exiting.")
		return
	}

	switch strings.TrimSpace(scanner.Text()) {
	case "1":
		if !runOSXDefaults(scanner, output) {
			return
		}
		fmt.Fprintln(output, "OSX defaults setup complete.")
	case "0":
		fmt.Fprintln(output, "Exiting.")
	default:
		fmt.Fprintln(output, "Invalid selection.")
	}
}

func runOSXDefaults(scanner *bufio.Scanner, output io.Writer) bool {
	currentComputerName, err := currentComputerName()
	if err != nil {
		fmt.Fprintf(output, "Failed to apply OSX defaults: %v\n", err)
		return false
	}

	computerName, err := promptComputerName(scanner, output, currentComputerName)
	if err != nil {
		fmt.Fprintf(output, "Failed to apply OSX defaults: %v\n", err)
		return false
	}

	if err := applyOSXDefaults(output, computerName); err != nil {
		fmt.Fprintf(output, "Failed to apply OSX defaults: %v\n", err)
		return false
	}
	return true
}

func currentComputerName() (string, error) {
	rawName, err := runCommandOutput("scutil", "--get", "ComputerName")
	if err != nil {
		return "", fmt.Errorf("failed to read current computer name: %w", err)
	}
	name := strings.TrimSpace(string(rawName))
	if name == "" {
		return "", fmt.Errorf("failed to read current computer name: empty value")
	}
	return name, nil
}

func promptComputerName(scanner *bufio.Scanner, output io.Writer, currentName string) (string, error) {
	fmt.Fprintf(output, "Current computer name: %s\n", currentName)
	fmt.Fprintln(output, "1) Keep existing name")
	fmt.Fprintln(output, "2) Replace with a new name")
	fmt.Fprint(output, "Select option (1 or 2): ")
	if !scanner.Scan() {
		return "", fmt.Errorf("selection is required")
	}
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		return currentName, nil
	case "2":
		fmt.Fprint(output, "Enter new computer name: ")
		if !scanner.Scan() {
			return "", fmt.Errorf("computer name is required")
		}
		computerName := strings.TrimSpace(scanner.Text())
		if computerName == "" {
			return "", fmt.Errorf("computer name is required")
		}
		return computerName, nil
	default:
		return "", fmt.Errorf("invalid selection: %s", choice)
	}
}

func applyOSXDefaults(output io.Writer, computerName string) error {
	home, err := userHomeDir()
	if err != nil {
		return err
	}

	screenshotFolder := filepath.Join(home, "Downloads")
	if err := os.MkdirAll(screenshotFolder, 0o755); err != nil {
		return fmt.Errorf("create screenshots folder: %w", err)
	}

	commands := application_settings.BuildMacOSDefaultsCommands(computerName, screenshotFolder)
	for _, command := range commands {
		if err := runCommand(command.Name, command.Args...); err != nil {
			if command.AllowFailure {
				fmt.Fprintf(output, "Ignoring failure for command: %s %s (%v)\n", command.Name, strings.Join(command.Args, " "), err)
				continue
			}
			return fmt.Errorf("%s %s: %w", command.Name, strings.Join(command.Args, " "), err)
		}
	}

	for _, app := range []string{"Address Book", "Calendar", "Contacts", "Dock", "Finder", "Mail", "Safari", "SystemUIServer", "iCal"} {
		if err := runCommand("killall", app); err != nil {
			fmt.Fprintf(output, "Ignoring failure for command: killall %s (%v)\n", app, err)
		}
	}

	return nil
}
