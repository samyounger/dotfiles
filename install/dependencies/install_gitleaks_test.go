package dependencies

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInstallGitleaks_AsdfMissing_Exits(t *testing.T) {
	execLookPath = func(file string) (string, error) {
		if file == "asdf" || file == "gitleaks" {
			return "", errors.New("not found")
		}
		return "/usr/local/bin/" + file, nil
	}
	defer resetExecMocks()
	defer func() { osExit = os.Exit }()
	osExit = func(code int) { panic("osExit called") }
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected osExit to be called if asdf is not installed")
		}
	}()
	InstallGitleaks()
}

func TestInstallGitleaks_PluginNotPresent(t *testing.T) {
	calls := [][]string{}
	pluginListCalled := false
	addCalled := false
	installCalled := false
	globalCalled := false
	versionCalled := false
	tempHome := t.TempDir()
	toolVersionsDir := filepath.Join(tempHome, ".dotfiles", "runcom")
	if err := os.MkdirAll(toolVersionsDir, 0o755); err != nil {
		t.Fatalf("mkdir tool versions dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(toolVersionsDir, ".tool-versions"), []byte("gitleaks 8.30.0\n"), 0o644); err != nil {
		t.Fatalf("write tool versions: %v", err)
	}
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("set HOME: %v", err)
	}
	defer os.Setenv("HOME", origHome)

	execLookPath = func(file string) (string, error) {
		if file == "gitleaks" {
			return "", errors.New("not found")
		}
		return "/usr/local/bin/" + file, nil
	}

	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		cmd := exec.Command("echo")
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
		if name == "gitleaks" && len(args) > 0 && args[0] == "version" {
			versionCalled = true
		}
		return cmd
	}
	defer resetExecMocks()

	InstallGitleaks()
	if !pluginListCalled || !addCalled || !installCalled || !globalCalled || !versionCalled {
		t.Errorf("Expected plugin-list, plugin-add, install, global, version; got calls: %v", calls)
	}
}

func TestInstallGitleaks_PluginPresent(t *testing.T) {
	calls := [][]string{}
	pluginListCalled := false
	addCalled := false
	installCalled := false
	globalCalled := false
	versionCalled := false
	tempHome := t.TempDir()
	toolVersionsDir := filepath.Join(tempHome, ".dotfiles", "runcom")
	if err := os.MkdirAll(toolVersionsDir, 0o755); err != nil {
		t.Fatalf("mkdir tool versions dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(toolVersionsDir, ".tool-versions"), []byte("gitleaks 8.30.0\n"), 0o644); err != nil {
		t.Fatalf("write tool versions: %v", err)
	}
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("set HOME: %v", err)
	}
	defer os.Setenv("HOME", origHome)

	execLookPath = func(file string) (string, error) {
		if file == "gitleaks" {
			return "", errors.New("not found")
		}
		return "/usr/local/bin/" + file, nil
	}

	execCommand = func(name string, args ...string) *exec.Cmd {
		calls = append(calls, append([]string{name}, args...))
		cmd := exec.Command("echo")
		if name == "asdf" && len(args) > 0 && args[0] == "plugin-list" {
			pluginListCalled = true
			cmd = exec.Command("echo", "gitleaks\n")
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
		if name == "gitleaks" && len(args) > 0 && args[0] == "version" {
			versionCalled = true
		}
		return cmd
	}
	defer resetExecMocks()

	InstallGitleaks()
	if !pluginListCalled || addCalled || !installCalled || !globalCalled || !versionCalled {
		t.Errorf("Expected plugin-list, install, global, version; plugin-add should not be called. Calls: %v", calls)
	}
}

func TestInstallGitleaks_CommandError(t *testing.T) {
	tempHome := t.TempDir()
	toolVersionsDir := filepath.Join(tempHome, ".dotfiles", "runcom")
	if err := os.MkdirAll(toolVersionsDir, 0o755); err != nil {
		t.Fatalf("mkdir tool versions dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(toolVersionsDir, ".tool-versions"), []byte("gitleaks 8.30.0\n"), 0o644); err != nil {
		t.Fatalf("write tool versions: %v", err)
	}
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("set HOME: %v", err)
	}
	defer os.Setenv("HOME", origHome)

	execLookPath = func(file string) (string, error) {
		if file == "gitleaks" {
			return "", errors.New("not found")
		}
		return "/usr/local/bin/" + file, nil
	}

	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "asdf" && len(args) > 0 && args[0] == "install" {
			return exec.Command("false")
		}
		return exec.Command("echo")
	}
	defer resetExecMocks()
	defer func() { osExit = os.Exit }()
	osExit = func(code int) { panic("osExit called") }
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected osExit to be called on command error")
		}
	}()

	InstallGitleaks()
}

func TestInstallGitleaks_AlreadyInstalled(t *testing.T) {
	called := false
	execLookPath = func(file string) (string, error) {
		if file == "gitleaks" {
			return "/usr/local/bin/gitleaks", nil
		}
		return origExecLookPath(file)
	}
	execCommand = func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo")
	}
	defer resetExecMocks()

	InstallGitleaks()
	if !called {
		t.Errorf("Expected execCommand to be called if gitleaks is already installed")
	}
}

func TestDesiredGitleaksVersion(t *testing.T) {
	tempHome := t.TempDir()
	toolVersionsDir := filepath.Join(tempHome, ".dotfiles", "runcom")
	if err := os.MkdirAll(toolVersionsDir, 0o755); err != nil {
		t.Fatalf("mkdir tool versions dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(toolVersionsDir, ".tool-versions"), []byte("golang 1.26.1\ngitleaks 8.30.0\n"), 0o644); err != nil {
		t.Fatalf("write tool versions: %v", err)
	}
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("set HOME: %v", err)
	}
	defer os.Setenv("HOME", origHome)

	version, err := desiredGitleaksVersion()
	if err != nil {
		t.Fatalf("expected version, got error: %v", err)
	}
	if version != "8.30.0" {
		t.Fatalf("expected 8.30.0, got %s", version)
	}
}
