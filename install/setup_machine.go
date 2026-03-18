package install

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/samyounger/dotfiles/install/application_settings"
)

var runCommandOutput = func(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}
var desiredDockApplicationsFunc = application_settings.DesiredDockApplications
var buildDockFindCommandFunc = application_settings.BuildDockFindCommand
var buildDockAddCommandFunc = application_settings.BuildDockAddCommand
var buildDockRestartCommandFunc = application_settings.BuildDockRestartCommand

func SetupMachine(input io.Reader, output io.Writer) {
	fmt.Fprintln(output, "Machine Setup Options:")
	fmt.Fprintln(output, "1) Setup MacOS settings")
	fmt.Fprintln(output, "2) Setup chrome")
	fmt.Fprintln(output, "3) Setup dock")
	fmt.Fprintln(output, "0) Exit")
	fmt.Fprint(output, "$ ")

	scanner := bufio.NewScanner(input)
	if !scanner.Scan() {
		fmt.Fprintln(output, "Exiting.")
		return
	}

	switch strings.TrimSpace(scanner.Text()) {
	case "1":
		if !runMacOSDefaults(scanner, output) {
			return
		}
		fmt.Fprintln(output, "MacOS settings setup complete.")
	case "2":
		if err := applyChromeDefaults(); err != nil {
			fmt.Fprintf(output, "Failed to apply Chrome defaults: %v\n", err)
			return
		}
		fmt.Fprintln(output, "Chrome defaults setup complete.")
	case "3":
		if err := applyDockDefaults(); err != nil {
			fmt.Fprintf(output, "Failed to apply Dock defaults: %v\n", err)
			return
		}
		fmt.Fprintln(output, "Dock defaults setup complete.")
	case "0":
		fmt.Fprintln(output, "Exiting.")
	default:
		fmt.Fprintln(output, "Invalid selection.")
	}
}

func runMacOSDefaults(scanner *bufio.Scanner, output io.Writer) bool {
	currentComputerName, err := currentComputerName()
	if err != nil {
		fmt.Fprintf(output, "Failed to apply MacOS settings: %v\n", err)
		return false
	}

	computerName, err := promptComputerName(scanner, output, currentComputerName)
	if err != nil {
		fmt.Fprintf(output, "Failed to apply MacOS settings: %v\n", err)
		return false
	}

	if err := applyMacOSDefaults(output, computerName); err != nil {
		fmt.Fprintf(output, "Failed to apply MacOS settings: %v\n", err)
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

func applyMacOSDefaults(output io.Writer, computerName string) error {
	home, err := userHomeDir()
	if err != nil {
		return err
	}

	screenshotFolder := filepath.Join(home, "Downloads")
	if err := os.MkdirAll(screenshotFolder, 0o755); err != nil {
		return fmt.Errorf("create screenshots folder: %w", err)
	}

	commands := application_settings.BuildMacOSSettingsCommands(computerName, screenshotFolder)
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

func applyChromeDefaults() error {
	commands := application_settings.BuildChromeDefaultsCommands()
	for _, command := range commands {
		if err := runCommand(command.Name, command.Args...); err != nil {
			return fmt.Errorf("%s %s: %w", command.Name, strings.Join(command.Args, " "), err)
		}
	}
	return nil
}

func applyDockDefaults() error {
	changed := false
	for _, appPath := range desiredDockApplicationsFunc() {
		exists, err := dockAppExists(appPath)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		command := buildDockAddCommandFunc(appPath)
		if err := runCommand(command.Name, command.Args...); err != nil {
			return fmt.Errorf("%s %s: %w", command.Name, strings.Join(command.Args, " "), err)
		}
		changed = true
	}

	if !changed {
		return nil
	}

	restartCommand := buildDockRestartCommandFunc()
	if err := runCommand(restartCommand.Name, restartCommand.Args...); err != nil {
		return fmt.Errorf("%s %s: %w", restartCommand.Name, strings.Join(restartCommand.Args, " "), err)
	}

	return nil
}

func dockAppExists(appPath string) (bool, error) {
	command := buildDockFindCommandFunc(appPath)
	_, err := runCommandOutput(command.Name, command.Args...)
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("%s %s: %w", command.Name, strings.Join(command.Args, " "), err)
}
