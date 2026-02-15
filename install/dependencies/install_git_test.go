package dependencies

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func mockExecLookPath(found bool) func(string) (string, error) {
	return func(file string) (string, error) {
		if found {
			return "/usr/local/bin/" + file, nil
		}
		return "", errors.New("not found")
	}
}

func TestInstallGit_AlreadyInstalled(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "git" {
			return "/usr/local/bin/git", nil
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer resetExecMocks()

	InstallGit()
	if !called {
		t.Errorf("Expected execCommand to be called to print git version when git is already installed")
	}
}

func TestCheckGitInstalled_Found(t *testing.T) {
	called := false
	execLookPath = mockExecLookPath(true)
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer resetExecMocks()
	checkGitInstalled()
	if !called {
		t.Errorf("Expected execCommand to be called for git --version")
	}
}

func TestCheckGitInstalled_NotFound(t *testing.T) {
	execLookPath = mockExecLookPath(false)
	defer resetExecMocks()
	checkGitInstalled()
}

func TestInstallGit_Success(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo")
	}
	defer resetExecMocks()
	installGit()
	if len(calls) != 2 || calls[0][0] != "brew" || calls[1][0] != "git" {
		t.Errorf("Expected brew install git and git --version, got %v", calls)
	}
}

func TestInstallGit_BrewError(t *testing.T) {
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
	installGit()
}
