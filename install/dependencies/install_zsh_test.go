package dependencies

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestInstallZsh_AlreadyInstalled(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "zsh" {
			return "/usr/local/bin/zsh", nil
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer resetExecMocks()

	InstallZsh()
	if !called {
		t.Errorf("Expected execCommand to be called to print zsh version")
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
	defer resetExecMocks()

	// Mock exec.LookPath to simulate zsh not installed
	execLookPath = func(file string) (string, error) {
		if file == "zsh" {
			return "", errors.New("not found")
		}
		return origExecLookPath(file)
	}

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
	defer resetExecMocks()

	// Mock exec.LookPath to simulate zsh not installed
	execLookPath = func(file string) (string, error) {
		if file == "zsh" {
			return "", errors.New("not found")
		}
		return origExecLookPath(file)
	}

	defer func() { osExit = os.Exit }()
	osExit = func(code int) { panic("osExit called") }
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected osExit to be called on brew install error")
		}
	}()
	InstallZsh()
}
