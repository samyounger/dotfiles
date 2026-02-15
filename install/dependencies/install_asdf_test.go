package dependencies

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestCheckAsdfInstalled_Found(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "asdf" {
			return "/usr/local/bin/asdf", nil
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer resetExecMocks()
	defer func() { osExit = os.Exit }()
	osExit = func(code int) { panic("osExit called") }
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected osExit to be called")
		}
	}()
	checkAsdfInstalled()
	if !called {
		t.Errorf("Expected execCommand to be called for asdf --version")
	}
}

func TestCheckAsdfInstalled_NotFound(t *testing.T) {
	execLookPath = func(file string) (string, error) {
		if file == "asdf" {
			return "", errors.New("not found")
		}
		return origExecLookPath(file)
	}
	defer resetExecMocks()
	checkAsdfInstalled()
}

func TestInstallAsdf_Success(t *testing.T) {
	calls := [][]string{}
	execLookPath = func(file string) (string, error) {
		if file == "asdf" {
			return "", errors.New("not found")
		}
		if file == "brew" {
			return "/usr/local/bin/brew", nil
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo")
	}
	defer resetExecMocks()

	installAsdf()
	if len(calls) != 2 || calls[0][0] != "brew" || calls[1][0] != "asdf" {
		t.Errorf("Expected brew install asdf and asdf --version, got %v", calls)
	}
}

func TestInstallAsdf_BrewError(t *testing.T) {
	execLookPath = func(file string) (string, error) {
		if file == "asdf" {
			return "", errors.New("not found")
		}
		if file == "brew" {
			return "/usr/local/bin/brew", nil
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer resetExecMocks()
	defer func() { osExit = os.Exit }()
	osExit = func(code int) { panic("osExit called") }
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected osExit to be called on error")
		}
	}()
	installAsdf()
}
