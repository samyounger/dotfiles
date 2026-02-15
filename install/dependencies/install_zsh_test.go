package dependencies

import (
	"os/exec"
	"testing"
	"errors"
)

var (
	origExecCommand = execCommand
)

func resetExecZshMocks() {
	execCommand = origExecCommand
}

func TestInstallZsh_AlreadyInstalled(t *testing.T) {
	called := false
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer resetExecZshMocks()

	// Mock exec.LookPath to simulate zsh already installed
	origLookPath := execLookPath
	execLookPath = func(file string) (string, error) {
		if file == "zsh" {
			return "/usr/local/bin/zsh", nil
		}
		return origLookPath(file)
	}
	defer func() { execLookPath = origLookPath }()

	InstallZsh()
	if called {
		t.Errorf("Did not expect execCommand to be called if zsh is already installed")
	}
}

func TestInstallZsh_NotInstalled_Success(t *testing.T) {
	called := false
	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "brew" && len(args) > 0 && args[0] == "install" && args[1] == "zsh" {
			called = true
		}
		return exec.Command("echo")
	}
	defer resetExecZshMocks()

	// Mock exec.LookPath to simulate zsh not installed
	origLookPath := execLookPath
	execLookPath = func(file string) (string, error) {
		if file == "zsh" {
			return "", errors.New("not found")
		}
		return origLookPath(file)
	}
	defer func() { execLookPath = origLookPath }()

	InstallZsh()
	if !called {
		t.Errorf("Expected execCommand to be called to install zsh")
	}
}

func TestInstallZsh_BrewError(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "brew" && len(args) > 0 && args[0] == "install" && args[1] == "zsh" {
			return exec.Command("false")
		}
		return exec.Command("echo")
	}
	defer resetExecZshMocks()

	// Mock exec.LookPath to simulate zsh not installed
	origLookPath := execLookPath
	execLookPath = func(file string) (string, error) {
		if file == "zsh" {
			return "", errors.New("not found")
		}
		return origLookPath(file)
	}
	defer func() { execLookPath = origLookPath }()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected os.Exit to be called on brew install error")
		}
	}()
	InstallZsh()
}

