package dependencies

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestInstallHomebrew_AlreadyInstalled(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "brew" {
			return "/usr/local/bin/brew", nil
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer resetExecMocks()

	InstallHomebrew()
	if !called {
		t.Errorf("Expected execCommand to be called to print Homebrew version")
	}
}

func TestInstallHomebrew_NotInstalled_Success(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "brew" {
			return "", errors.New("not found")
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "/bin/bash" && len(args) > 1 && args[0] == "-c" {
			called = true
		}
		return exec.Command("echo")
	}
	defer resetExecMocks()

	InstallHomebrew()
	if !called {
		t.Errorf("Expected execCommand to be called to install Homebrew")
	}
}

func TestInstallHomebrew_InstallError(t *testing.T) {
	execLookPath = func(file string) (string, error) {
		if file == "brew" {
			return "", errors.New("not found")
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "/bin/bash" && len(args) > 1 && args[0] == "-c" {
			return exec.Command("false")
		}
		return exec.Command("echo")
	}
	defer resetExecMocks()
	defer func() { osExit = os.Exit }()
	osExit = func(code int) { panic("osExit called") }
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected osExit to be called on Homebrew install error")
		}
	}()
	InstallHomebrew()
}
