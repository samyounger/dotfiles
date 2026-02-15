package install

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/syounger-ut/dotfiles/install/application_settings"
)

type commandCall struct {
	name string
	args []string
}

func TestSetupMachine_Exit(t *testing.T) {
	defer resetOSXDefaultsMocks()

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("0\n"), out)

	if !strings.Contains(out.String(), "Exiting.") {
		t.Fatalf("expected exit message, got: %s", out.String())
	}
}

func TestSetupMachine_InvalidTopLevelSelection(t *testing.T) {
	defer resetOSXDefaultsMocks()

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("9\n"), out)

	if !strings.Contains(out.String(), "Invalid selection.") {
		t.Fatalf("expected invalid selection message, got: %s", out.String())
	}
}

func TestSetupMachine_SetupChromeSuccess(t *testing.T) {
	defer resetOSXDefaultsMocks()

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("2\n"), out)

	if !strings.Contains(out.String(), "Chrome defaults setup complete.") {
		t.Fatalf("expected chrome success message, got: %s", out.String())
	}
	if !hasCommand(calls, "defaults", "write", "com.google.Chrome", "DisablePrintPreview", "-bool", "true") {
		t.Fatalf("expected chrome defaults command to be run, got %#v", calls)
	}
}

func TestSetupMachine_SetupChromeFailure(t *testing.T) {
	defer resetOSXDefaultsMocks()

	runCommand = func(name string, args ...string) error {
		if name == "defaults" {
			return errors.New("defaults failed")
		}
		return nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("2\n"), out)

	if !strings.Contains(out.String(), "Failed to apply Chrome defaults") {
		t.Fatalf("expected chrome failure message, got: %s", out.String())
	}
}

func TestSetupMachine_SetupDockSuccess(t *testing.T) {
	defer resetOSXDefaultsMocks()

	desiredDockApplicationsFunc = func() []string { return []string{"app1", "app2"} }
	buildDockFindCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "find", Args: []string{appPath}}
	}
	buildDockAddCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "add", Args: []string{appPath}}
	}
	buildDockRestartCommandFunc = func() application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "restart", Args: []string{"Dock"}}
	}

	notFoundErr := exitStatusError(t, 1)
	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		if name == "find" {
			return nil, notFoundErr
		}
		return nil, errors.New("unexpected command")
	}

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("3\n"), out)

	if !strings.Contains(out.String(), "Dock defaults setup complete.") {
		t.Fatalf("expected dock success message, got: %s", out.String())
	}
	if !hasCommand(calls, "add", "app1") {
		t.Fatalf("expected dock command to be run, got %#v", calls)
	}
	if !hasCommand(calls, "add", "app2") {
		t.Fatalf("expected second dock add command to be run, got %#v", calls)
	}
	if !hasCommand(calls, "restart", "Dock") {
		t.Fatalf("expected Dock restart command to be run, got %#v", calls)
	}
}

func TestSetupMachine_SetupDockFailure(t *testing.T) {
	defer resetOSXDefaultsMocks()

	desiredDockApplicationsFunc = func() []string { return []string{"app1"} }
	buildDockFindCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "find", Args: []string{appPath}}
	}
	buildDockAddCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "add", Args: []string{appPath}}
	}

	notFoundErr := exitStatusError(t, 1)
	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		if name == "find" {
			return nil, notFoundErr
		}
		return nil, errors.New("unexpected command")
	}

	runCommand = func(name string, args ...string) error {
		if name == "add" {
			return errors.New("add failed")
		}
		return nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("3\n"), out)

	if !strings.Contains(out.String(), "Failed to apply Dock defaults") {
		t.Fatalf("expected dock failure message, got: %s", out.String())
	}
}

func TestSetupMachine_SuccessKeepExistingName(t *testing.T) {
	defer resetOSXDefaultsMocks()

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }
	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		if name == "scutil" && len(args) == 2 && args[0] == "--get" && args[1] == "ComputerName" {
			return []byte("CurrentMac\n"), nil
		}
		return nil, errors.New("unexpected command")
	}
	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("1\n1\n"), out)

	if !strings.Contains(out.String(), "Current computer name: CurrentMac") {
		t.Fatalf("expected current name prompt, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "MacOS settings setup complete.") {
		t.Fatalf("expected success message, got: %s", out.String())
	}
	if !hasCommand(calls, "sudo", "scutil", "--set", "ComputerName", "CurrentMac") {
		t.Fatalf("expected current computer name command, got %#v", calls)
	}
}

func TestSetupMachine_SuccessReplaceName(t *testing.T) {
	defer resetOSXDefaultsMocks()

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }
	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		return []byte("CurrentMac\n"), nil
	}

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("1\n2\nWorkstation\n"), out)

	if !hasCommand(calls, "sudo", "scutil", "--set", "ComputerName", "Workstation") {
		t.Fatalf("expected replaced computer name command, got %#v", calls)
	}
}

func TestSetupMachine_FailsWhenCurrentNameLookupFails(t *testing.T) {
	defer resetOSXDefaultsMocks()

	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("lookup failed")
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("1\n"), out)

	if !strings.Contains(out.String(), "Failed to apply MacOS settings") {
		t.Fatalf("expected failure message, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "failed to read current computer name") {
		t.Fatalf("expected current-name lookup error, got: %s", out.String())
	}
}

func TestSetupMachine_RequiresValidNameSelection(t *testing.T) {
	defer resetOSXDefaultsMocks()

	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		return []byte("CurrentMac\n"), nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("1\n3\n"), out)

	if !strings.Contains(out.String(), "invalid selection") {
		t.Fatalf("expected invalid selection error, got: %s", out.String())
	}
}

func TestSetupMachine_RequiresComputerNameWhenReplacing(t *testing.T) {
	defer resetOSXDefaultsMocks()

	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		return []byte("CurrentMac\n"), nil
	}

	out := &bytes.Buffer{}
	SetupMachine(strings.NewReader("1\n2\n\n"), out)

	if !strings.Contains(out.String(), "computer name is required") {
		t.Fatalf("expected required-name error, got: %s", out.String())
	}
}

func TestApplyOSXDefaults_UsesInputComputerName(t *testing.T) {
	defer resetOSXDefaultsMocks()

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	out := &bytes.Buffer{}
	if err := applyMacOSDefaults(out, "Workstation"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasCommand(calls, "sudo", "scutil", "--set", "ComputerName", "Workstation") {
		t.Fatalf("expected input computer name command, got %#v", calls)
	}

	screenshotFolder := filepath.Join(home, "Downloads")
	if !hasCommand(calls, "defaults", "write", "com.apple.screencapture", "location", "-string", screenshotFolder) {
		t.Fatalf("expected screenshot location command for %s", screenshotFolder)
	}
}

func TestApplyOSXDefaults_AllowFailureContinues(t *testing.T) {
	defer resetOSXDefaultsMocks()

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	runCommand = func(name string, args ...string) error {
		if name == "osascript" {
			return errors.New("not available")
		}
		return nil
	}

	out := &bytes.Buffer{}
	if err := applyMacOSDefaults(out, "Workstation"); err != nil {
		t.Fatalf("expected allowFailure command to be ignored, got error: %v", err)
	}
	if !strings.Contains(out.String(), "Ignoring failure for command: osascript") {
		t.Fatalf("expected allowFailure log message, got: %s", out.String())
	}
}

func TestApplyOSXDefaults_NonAllowFailureReturnsError(t *testing.T) {
	defer resetOSXDefaultsMocks()

	home := t.TempDir()
	userHomeDir = func() (string, error) { return home, nil }

	runCommand = func(name string, args ...string) error {
		if name == "sudo" && len(args) == 1 && args[0] == "-v" {
			return errors.New("auth failed")
		}
		return nil
	}

	out := &bytes.Buffer{}
	err := applyMacOSDefaults(out, "Workstation")
	if err == nil {
		t.Fatalf("expected error from non-allowFailure command")
	}
	if !strings.Contains(err.Error(), "sudo -v") {
		t.Fatalf("expected command context in error, got: %v", err)
	}
}

func TestApplyChromeDefaults_Success(t *testing.T) {
	defer resetOSXDefaultsMocks()

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	if err := applyChromeDefaults(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(calls) != 6 {
		t.Fatalf("expected 6 chrome defaults commands, got %d", len(calls))
	}
}

func TestApplyChromeDefaults_Failure(t *testing.T) {
	defer resetOSXDefaultsMocks()

	runCommand = func(name string, args ...string) error {
		if name == "defaults" && len(args) > 0 && args[0] == "write" {
			return errors.New("boom")
		}
		return nil
	}

	err := applyChromeDefaults()
	if err == nil {
		t.Fatalf("expected error from chrome defaults command")
	}
	if !strings.Contains(err.Error(), "defaults write") {
		t.Fatalf("expected contextual command error, got: %v", err)
	}
}

func TestApplyDockDefaults_Success(t *testing.T) {
	defer resetOSXDefaultsMocks()

	desiredDockApplicationsFunc = func() []string { return []string{"app1", "app2"} }
	buildDockFindCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "find", Args: []string{appPath}}
	}
	buildDockAddCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "add", Args: []string{appPath}}
	}
	buildDockRestartCommandFunc = func() application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "restart", Args: []string{"Dock"}}
	}

	notFoundErr := exitStatusError(t, 1)
	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		if name == "find" {
			return nil, notFoundErr
		}
		return nil, errors.New("unexpected command")
	}

	var calls []commandCall
	runCommand = func(name string, args ...string) error {
		calls = append(calls, commandCall{name: name, args: append([]string{}, args...)})
		return nil
	}

	if err := applyDockDefaults(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(calls) != 3 {
		t.Fatalf("expected 3 dock execution commands, got %d", len(calls))
	}
}

func TestApplyDockDefaults_Failure(t *testing.T) {
	defer resetOSXDefaultsMocks()

	desiredDockApplicationsFunc = func() []string { return []string{"app1"} }
	buildDockFindCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "find", Args: []string{appPath}}
	}
	buildDockAddCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "add", Args: []string{appPath}}
	}

	notFoundErr := exitStatusError(t, 1)
	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		if name == "find" {
			return nil, notFoundErr
		}
		return nil, errors.New("unexpected command")
	}

	runCommand = func(name string, args ...string) error {
		if name == "add" {
			return errors.New("boom")
		}
		return nil
	}

	err := applyDockDefaults()
	if err == nil {
		t.Fatalf("expected error from dock defaults command")
	}
	if !strings.Contains(err.Error(), "add app1") {
		t.Fatalf("expected contextual command error, got: %v", err)
	}
}

func TestApplyDockDefaults_NoChangesWhenAppsAlreadyExist(t *testing.T) {
	defer resetOSXDefaultsMocks()

	desiredDockApplicationsFunc = func() []string { return []string{"app1", "app2"} }
	buildDockFindCommandFunc = func(appPath string) application_settings.SettingsCommand {
		return application_settings.SettingsCommand{Name: "find", Args: []string{appPath}}
	}

	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		if name == "find" {
			return []byte("found"), nil
		}
		return nil, errors.New("unexpected command")
	}

	called := false
	runCommand = func(name string, args ...string) error {
		called = true
		return nil
	}

	if err := applyDockDefaults(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatalf("did not expect dock mutation commands when all apps already exist")
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

func resetOSXDefaultsMocks() {
	userHomeDir = os.UserHomeDir
	runCommand = func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	runCommandOutput = func(name string, args ...string) ([]byte, error) {
		cmd := exec.Command(name, args...)
		cmd.Stderr = os.Stderr
		return cmd.Output()
	}
	desiredDockApplicationsFunc = application_settings.DesiredDockApplications
	buildDockFindCommandFunc = application_settings.BuildDockFindCommand
	buildDockAddCommandFunc = application_settings.BuildDockAddCommand
	buildDockRestartCommandFunc = application_settings.BuildDockRestartCommand
}

func exitStatusError(t *testing.T, code int) error {
	t.Helper()
	_, err := exec.Command("sh", "-c", fmt.Sprintf("exit %d", code)).Output()
	if err == nil {
		t.Fatalf("expected non-nil error for exit code %d", code)
	}
	return err
}
