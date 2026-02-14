package dependencies

import (
	"errors"
	"os/exec"
	"testing"
)

// Mocks for exec.LookPath and exec.Command
var (
	origExecLookPath = execLookPath
	origExecCommand = execCommand
)

func mockExecLookPath(found bool) func(string) (string, error) {
	return func(file string) (string, error) {
		if found {
			return "/usr/local/bin/git", nil
		}
		return "", errors.New("not found")
	}
}

func mockExecCommand(runErr error) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cmd := exec.Command("echo")
		cmd.Run = func() error { return runErr }
		return cmd
	}
}

func TestCheckGitInstalled_Found(t *testing.T) {
	called := false
	execLookPath = mockExecLookPath(true)
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer func() {
		execLookPath = origExecLookPath
		execCommand = origExecCommand
		if r := recover(); r != nil {
			// os.Exit will panic in test, that's expected
		}
	}()
	checkGitInstalled()
	if !called {
		t.Errorf("Expected execCommand to be called for git --version")
	}
}

func TestCheckGitInstalled_NotFound(t *testing.T) {
	execLookPath = mockExecLookPath(false)
	defer func() { execLookPath = origExecLookPath }()
	checkGitInstalled()
}

func TestInstallGit_Success(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo")
	}
	defer func() { execCommand = origExecCommand }()
	installGit()
	if len(calls) != 2 || calls[0][0] != "brew" || calls[1][0] != "git" {
		t.Errorf("Expected brew install git and git --version, got %v", calls)
	}
}

func TestInstallGit_BrewError(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	defer func() {
		execCommand = origExecCommand
		if r := recover(); r == nil {
			t.Errorf("Expected os.Exit to be called on error")
		}
	}()
	installGit()
}

