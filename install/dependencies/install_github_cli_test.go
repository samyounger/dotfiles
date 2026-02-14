package dependencies

import (
	"errors"
	"os/exec"
	"testing"
)

var (
	origExecLookPath = execLookPath
	origExecCommand   = execCommand
)

func resetExecMocks() {
	execLookPath = origExecLookPath
	execCommand = origExecCommand
}

func TestInstallGithubCLI_AsdfNotInstalled(t *testing.T) {
	execLookPath = func(file string) (string, error) {
		if file == "asdf" {
			return "", errors.New("not found")
		}
		return "/usr/local/bin/" + file, nil
	}
	defer resetExecMocks()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected os.Exit to be called if asdf is not installed")
		}
	}()
	InstallGithubCLI()
}

func TestInstallGithubCLI_PluginNotPresent(t *testing.T) {
	calls := [][]string{}
	pluginListCalled := false
	addCalled := false
	installCalled := false
	globalCalled := false
	versionCalled := false

	execLookPath = func(file string) (string, error) {
		return "/usr/local/bin/" + file, nil
	}

	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		cmd := exec.Command("echo")
		// Simulate plugin-list output (no gh plugin)
		if name == "asdf" && len(args) > 0 && args[0] == "plugin-list" {
			pluginListCalled = true
			cmd = exec.Command("echo", "")
		}
		if name == "asdf" && len(args) > 0 && args[0] == "plugin-add" {
			addCalled = true
		}
		if name == "asdf" && len(args) > 0 && args[0] == "install" {
			installCalled = true
		}
		if name == "asdf" && len(args) > 0 && args[0] == "global" {
			globalCalled = true
		}
		if name == "gh" && len(args) > 0 && args[0] == "--version" {
			versionCalled = true
		}
		return cmd
	}
	defer resetExecMocks()
	InstallGithubCLI()
	if !pluginListCalled || !addCalled || !installCalled || !globalCalled || !versionCalled {
		t.Errorf("Expected all asdf/gh commands to be called, got calls: %v", calls)
	}
}

func TestInstallGithubCLI_PluginPresent(t *testing.T) {
	calls := [][]string{}
	pluginListCalled := false
	addCalled := false
	installCalled := false
	globalCalled := false
	versionCalled := false

	execLookPath = func(file string) (string, error) {
		return "/usr/local/bin/" + file, nil
	}

	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		cmd := exec.Command("echo")
		// Simulate plugin-list output (gh plugin present)
		if name == "asdf" && len(args) > 0 && args[0] == "plugin-list" {
			pluginListCalled = true
			cmd = exec.Command("echo", "gh\n")
		}
		if name == "asdf" && len(args) > 0 && args[0] == "plugin-add" {
			addCalled = true
		}
		if name == "asdf" && len(args) > 0 && args[0] == "install" {
			installCalled = true
		}
		if name == "asdf" && len(args) > 0 && args[0] == "global" {
			globalCalled = true
		}
		if name == "gh" && len(args) > 0 && args[0] == "--version" {
			versionCalled = true
		}
		return cmd
	}
	defer resetExecMocks()
	InstallGithubCLI()
	if !pluginListCalled || addCalled || !installCalled || !globalCalled || !versionCalled {
		t.Errorf("Expected plugin-list, install, global, version; plugin-add should not be called. Calls: %v", calls)
	}
}

func TestInstallGithubCLI_CommandError(t *testing.T) {
	calls := [][]string{}
	failOn := "install"

	execLookPath = func(file string) (string, error) {
		return "/usr/local/bin/" + file, nil
	}

	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		if name == "asdf" && len(args) > 0 && args[0] == failOn {
			return exec.Command("false")
		}
		return exec.Command("echo")
	}
	defer resetExecMocks()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected os.Exit to be called on command error")
		}
	}()
	InstallGithubCLI()
}

func TestInstallGithubCLI_PluginAddSSHFailHTTPSOk(t *testing.T) {
	calls := [][]string{}
	pluginListCalled := false
	addSSHCalled := false
	addHTTPSCalled := false
	addHTTPSEnv := []string{}
	installCalled := false
	globalCalled := false
	versionCalled := false

	execLookPath = func(file string) (string, error) {
		return "/usr/local/bin/" + file, nil
	}

	sshFail := false

	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		cmd := exec.Command("echo")
		// Simulate plugin-list output (no gh plugin)
		if name == "asdf" && len(args) > 0 && args[0] == "plugin-list" {
			pluginListCalled = true
			cmd = exec.Command("echo", "")
		}
		if name == "asdf" && len(args) > 0 && args[0] == "plugin-add" {
			if len(args) > 2 && args[2] == "git@github.com:kennyp/asdf-gh.git" {
				addSSHCalled = true
				sshFail = true
				// Simulate SSH failure by returning a command that fails
				return exec.Command("false")
			}
			if len(args) > 2 && args[2] == "https://github.com/kennyp/asdf-gh.git" {
				addHTTPSCalled = true
				addHTTPSEnv = os.Environ()
				cmd = exec.Command("echo")
			}
		}
		if name == "asdf" && len(args) > 0 && args[0] == "install" {
			installCalled = true
		}
		if name == "asdf" && len(args) > 0 && args[0] == "global" {
			globalCalled = true
		}
		if name == "gh" && len(args) > 0 && args[0] == "--version" {
			versionCalled = true
		}
		return cmd
	}
	defer resetExecMocks()
	InstallGithubCLI()
	if !pluginListCalled || !addSSHCalled || !addHTTPSCalled || !installCalled || !globalCalled || !versionCalled {
		t.Errorf("Expected SSH and HTTPS plugin-add, install, global, version; got calls: %v", calls)
	}
	// Check that GIT_ASKPASS is set in the environment for HTTPS
	askpassSet := false
	for _, env := range addHTTPSEnv {
		if env == "GIT_ASKPASS=/bin/echo" {
			askpassSet = true
			break
		}
	}
	if !askpassSet {
		t.Errorf("Expected GIT_ASKPASS to be set in HTTPS plugin-add environment")
	}
}

func TestInstallGithubCLI_AlreadyInstalled(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "gh" {
			return "/usr/local/bin/gh", nil
		}
		return "/usr/local/bin/" + file, nil
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "gh" && len(args) > 0 && args[0] == "--version" {
			called = true
		}
		return exec.Command("echo")
	}
	defer resetExecMocks()
	InstallGithubCLI()
	if !called {
		t.Errorf("Expected printGithubCLIVersion to be called if gh is already installed")
	}
}
