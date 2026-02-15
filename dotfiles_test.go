package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestMain_InstallGit(t *testing.T) {
	var calledGit, calledAsdf bool
	installGitFunc = func() { calledGit = true }
	installAsdfFunc = func() { calledAsdf = true }
	inputReader = strings.NewReader("1\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !calledGit {
		t.Errorf("Expected installGitFunc to be called")
	}
	if calledAsdf {
		t.Errorf("Did not expect installAsdfFunc to be called")
	}
	if !strings.Contains(buf.String(), "Options:") {
		t.Errorf("Expected menu to be printed")
	}
}

func TestMain_InstallAsdf(t *testing.T) {
	var calledGit, calledAsdf bool
	installGitFunc = func() { calledGit = true }
	installAsdfFunc = func() { calledAsdf = true }
	inputReader = strings.NewReader("2\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !calledAsdf {
		t.Errorf("Expected installAsdfFunc to be called")
	}
	if calledGit {
		t.Errorf("Did not expect installGitFunc to be called")
	}
}

func TestMain_Exit(t *testing.T) {
	installGitFunc = func() { t.Errorf("Should not call installGitFunc") }
	installAsdfFunc = func() { t.Errorf("Should not call installAsdfFunc") }
	installGithubCLIFunc = func() { t.Errorf("Should not call installGithubCLIFunc") }
	setupFunc = func() { t.Errorf("Should not call setupFunc") }
	inputReader = strings.NewReader("5\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !strings.Contains(buf.String(), "Exiting.") {
		t.Errorf("Expected 'Exiting.' to be printed")
	}
}

func TestMain_Default(t *testing.T) {
	installGitFunc = func() { t.Errorf("Should not call installGitFunc") }
	installAsdfFunc = func() { t.Errorf("Should not call installAsdfFunc") }
	inputReader = strings.NewReader("foobar\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !strings.Contains(buf.String(), "Selected: foobar") {
		t.Errorf("Expected default handler to print input")
	}
}
