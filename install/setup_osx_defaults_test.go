package install

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type commandCall struct {
	name string
	args []string
}

func TestSetupOSXDefaults_Success(t *testing.T) {
	defer resetDotfilesMocks()
	setComputerNameEnv(t, "")

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	runCommand = func(name string, args ...string) error { return nil }

	out := &bytes.Buffer{}
	SetupOSXDefaults(out)

	if !strings.Contains(out.String(), "OSX defaults setup complete.") {
		t.Fatalf("expected success message, got: %s", out.String())
	}
}

func TestSetupOSXDefaults_Failure(t *testing.T) {
	defer resetDotfilesMocks()
	setComputerNameEnv(t, "")

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	runCommand = func(name string, args ...string) error {
		if name == "sudo" && len(args) == 1 && args[0] == "-v" {
			return errors.New("boom")
		}
		return nil
	}

	out := &bytes.Buffer{}
	SetupOSXDefaults(out)

	if !strings.Contains(out.String(), "Failed to apply OSX defaults") {
		t.Fatalf("expected osxdefaults error message, got: %s", out.String())
	}
}

func TestApplyOSXDefaults_UsesDefaultComputerName(t *testing.T) {
	defer resetDotfilesMocks()
	setComputerNameEnv(t, "")

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	out := &bytes.Buffer{}
	if err := applyOSXDefaults(out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasCommand(calls, "sudo", "scutil", "--set", "ComputerName", defaultComputerName) {
		t.Fatalf("expected default computer name command, got %#v", calls)
	}

	screenshotFolder := filepath.Join(home, "Downloads")
	if !hasCommand(calls, "defaults", "write", "com.apple.screencapture", "location", "-string", screenshotFolder) {
		t.Fatalf("expected screenshot location command for %s", screenshotFolder)
	}

	if _, err := os.Stat(screenshotFolder); err != nil {
		t.Fatalf("expected screenshot folder to exist, got error: %v", err)
	}
}

func TestApplyOSXDefaults_UsesEnvComputerName(t *testing.T) {
	defer resetDotfilesMocks()
	setComputerNameEnv(t, "Workstation")

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	out := &bytes.Buffer{}
	if err := applyOSXDefaults(out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasCommand(calls, "sudo", "scutil", "--set", "ComputerName", "Workstation") {
		t.Fatalf("expected env computer name command, got %#v", calls)
	}
}

func TestApplyOSXDefaults_AllowFailureContinues(t *testing.T) {
	defer resetDotfilesMocks()
	setComputerNameEnv(t, "")

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	runCommand = func(name string, args ...string) error {
		if name == "osascript" {
			return errors.New("not available")
		}
		return nil
	}

	out := &bytes.Buffer{}
	if err := applyOSXDefaults(out); err != nil {
		t.Fatalf("expected allowFailure command to be ignored, got error: %v", err)
	}
	if !strings.Contains(out.String(), "Ignoring failure for command: osascript") {
		t.Fatalf("expected allowFailure log message, got: %s", out.String())
	}
}

func TestApplyOSXDefaults_NonAllowFailureReturnsError(t *testing.T) {
	defer resetDotfilesMocks()
	setComputerNameEnv(t, "")

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	runCommand = func(name string, args ...string) error {
		if name == "sudo" && len(args) == 1 && args[0] == "-v" {
			return errors.New("auth failed")
		}
		return nil
	}

	out := &bytes.Buffer{}
	err := applyOSXDefaults(out)
	if err == nil {
		t.Fatalf("expected error from non-allowFailure command")
	}
	if !strings.Contains(err.Error(), "sudo -v") {
		t.Fatalf("expected command context in error, got: %v", err)
	}
}

func TestBuildOSXDefaultsCommands_ContainsCoreCommands(t *testing.T) {
	commands := buildOSXDefaultsCommands("NeoTest", "/tmp/screens")
	if len(commands) == 0 {
		t.Fatal("expected command list to be non-empty")
	}

	if !hasOSXDefaultCommand(commands, "sudo", "scutil", "--set", "ComputerName", "NeoTest") {
		t.Fatalf("expected computer name command in command list")
	}

	if !hasOSXDefaultCommand(commands, "defaults", "write", "com.apple.screencapture", "location", "-string", "/tmp/screens") {
		t.Fatalf("expected screenshot location command in command list")
	}
}

func hasCommand(calls []commandCall, name string, args ...string) bool {
	for _, call := range calls {
		if call.name != name {
			continue
		}
		if len(call.args) != len(args) {
			continue
		}
		matched := true
		for i := range args {
			if call.args[i] != args[i] {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

func hasOSXDefaultCommand(commands []osxDefaultsCommand, name string, args ...string) bool {
	for _, command := range commands {
		if command.name != name {
			continue
		}
		if len(command.args) != len(args) {
			continue
		}
		matched := true
		for i := range args {
			if command.args[i] != args[i] {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

func setComputerNameEnv(t *testing.T, value string) {
	t.Helper()
	original, hadOriginal := os.LookupEnv("COMPUTER_NAME")
	if value == "" {
		if err := os.Unsetenv("COMPUTER_NAME"); err != nil {
			t.Fatalf("failed to unset COMPUTER_NAME: %v", err)
		}
	} else {
		if err := os.Setenv("COMPUTER_NAME", value); err != nil {
			t.Fatalf("failed to set COMPUTER_NAME: %v", err)
		}
	}
	t.Cleanup(func() {
		if hadOriginal {
			_ = os.Setenv("COMPUTER_NAME", original)
		} else {
			_ = os.Unsetenv("COMPUTER_NAME")
		}
	})
}
