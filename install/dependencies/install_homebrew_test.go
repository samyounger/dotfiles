package dependencies

import (
	"errors"
	"os/exec"
	"testing"
)

func TestInstallHomebrewIfNeeded_AlreadyInstalled(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "brew" {
			called = true
			return "/usr/local/bin/brew", nil
		}
		return "", errors.New("not found")
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		t.Errorf("execCommand should not be called when Homebrew is installed")
		return exec.Command("echo")
	}
	defer func() {
		execLookPath = exec.LookPath
		execCommand = exec.Command
	}()
	InstallHomebrewIfNeeded()
	if !called {
		t.Errorf("execLookPath was not called for brew")
	}
}

func TestInstallHomebrewIfNeeded_NotInstalled_Success(t *testing.T) {
	looked := false
	cmdCalled := false
	execLookPath = func(file string) (string, error) {
		if file == "brew" {
			looked = true
			return "", errors.New("not found")
		}
		return "", errors.New("not found")
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		cmdCalled = true
		return exec.Command("echo")
	}
	defer func() {
		execLookPath = exec.LookPath
		execCommand = exec.Command
	}()
	InstallHomebrewIfNeeded()
	if !looked {
		t.Errorf("execLookPath was not called for brew")
	}
	if !cmdCalled {
		t.Errorf("execCommand was not called to install Homebrew")
	}
}

func TestInstallHomebrewIfNeeded_NotInstalled_Fail(t *testing.T) {
	looked := false
	execLookPath = func(file string) (string, error) {
		if file == "brew" {
			looked = true
			return "", errors.New("not found")
		}
		return "", errors.New("not found")
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() {
		execLookPath = exec.LookPath
		execCommand = exec.Command
		if r := recover(); r == nil {
			t.Errorf("Expected os.Exit to be called on Homebrew install failure")
		}
	}()
	InstallHomebrewIfNeeded()
	if !looked {
		t.Errorf("execLookPath was not called for brew")
	}
}

