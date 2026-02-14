package dependencies

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

// mockCmd is a mock exec.Cmd
type mockCmd struct {
	args   []string
	runErr error
}

func (m *mockCmd) Run() error {
	return m.runErr
}

func (m *mockCmd) Output() ([]byte, error) {
	return []byte("asdf version 1.0.0"), nil
}

func (m *mockCmd) Start() error { return nil }
func (m *mockCmd) Wait() error  { return nil }
func (m *mockCmd) StdoutPipe() (os.ReadCloser, error) { return nil, nil }
func (m *mockCmd) StderrPipe() (os.ReadCloser, error) { return nil, nil }

// mockExecCommand returns a mockCmd
func mockExecCommand(expectedArgs ...string) func(name string, args ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}
}

// mockExecLookPath returns nil if asdf is found, error otherwise
func mockExecLookPath(found bool) func(string) (string, error) {
	return func(file string) (string, error) {
		if found {
			return "/usr/local/bin/asdf", nil
		}
		return "", errors.New("not found")
	}
}

func TestCheckAsdfInstalled_Found(t *testing.T) {
	called := false
	execLookPath = mockExecLookPath(true)
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer func() {
		execLookPath = exec.LookPath
		execCommand = exec.Command
		if r := recover(); r != nil {
			// os.Exit will panic in test, that's expected
		}
	}()
	checkAsdfInstalled()
	if !called {
		t.Errorf("Expected execCommand to be called for asdf --version")
	}
}

func TestCheckAsdfInstalled_NotFound(t *testing.T) {
	execLookPath = mockExecLookPath(false)
	defer func() { execLookPath = exec.LookPath }()
	// Should not exit or call execCommand
	checkAsdfInstalled()
}

func TestInstallAsdf_Success(t *testing.T) {
	calls := [][]string{}
	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		return exec.Command("echo")
	}
	defer func() { execCommand = exec.Command }()
	installAsdf()
	if len(calls) != 2 || calls[0][0] != "brew" || calls[1][0] != "asdf" {
		t.Errorf("Expected brew install asdf and asdf --version, got %v", calls)
	}
}

func TestInstallAsdf_BrewError(t *testing.T) {
	execCommand = func(name string, args ...string) *exec.Cmd {
		cmd := exec.Command("false")
		return cmd
	}
	defer func() {
		execCommand = exec.Command
		if r := recover(); r == nil {
			t.Errorf("Expected os.Exit to be called on error")
		}
	}()
	installAsdf()
}

