package install

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestSetup_All(t *testing.T) {
	var called []string
	options := []string{"Homebrew", "zsh", "git", "GitHub CLI", "asdf"}
	mockActions := make([]func(), len(options))
	for i, name := range options {
		name := name
		mockActions[i] = func() { called = append(called, name) }
	}

	// Patch os.Stdin and os.Stdout
	input := bytes.NewBufferString("0\n")
	output := &bytes.Buffer{}
	oldStdin := stdin
	oldStdout := stdout
	stdin = input
	stdout = output
	defer func() { stdin = oldStdin; stdout = oldStdout }()

	// Patch options in Setup
	origSetup := setupOptions
	setupOptions = func() []struct{ name string; action func() } {
		opts := make([]struct{ name string; action func() }, len(options))
		for i := range options {
			opts[i] = struct{ name string; action func() }{options[i], mockActions[i]}
		}
		return opts
	}
	defer func() { setupOptions = origSetup }()

	Setup()

	if len(called) != len(options) {
		t.Errorf("Expected all actions to be called, got %v", called)
	}
}

func TestSetup_SingleAndMultiple(t *testing.T) {
	var called []string
	options := []string{"Homebrew", "zsh", "git", "GitHub CLI", "asdf"}
	mockActions := make([]func(), len(options))
	for i, name := range options {
		name := name
		mockActions[i] = func() { called = append(called, name) }
	}

	cases := []struct {
		input    string
		expected []string
	}{
		{"1\n", []string{"Homebrew"}},
		{"2\n", []string{"zsh"}},
		{"1,3\n", []string{"Homebrew", "git"}},
		{"3,5\n", []string{"git", "asdf"}},
	}

	for _, tc := range cases {
		called = nil
		input := bytes.NewBufferString(tc.input)
		output := &bytes.Buffer{}
		oldStdin := stdin
		oldStdout := stdout
		stdin = input
		stdout = output

		origSetup := setupOptions
		setupOptions = func() []struct{ name string; action func() } {
			opts := make([]struct{ name string; action func() }, len(options))
			for i := range options {
				opts[i] = struct{ name string; action func() }{options[i], mockActions[i]}
			}
			return opts
		}

		Setup()

		if len(called) != len(tc.expected) {
			fmt.Printf("input: %q\n", tc.input)
			t.Errorf("Expected %v, got %v", tc.expected, called)
		}
		for i, v := range tc.expected {
			if called[i] != v {
				t.Errorf("Expected %v, got %v", tc.expected, called)
			}
		}

		stdin = oldStdin
		stdout = oldStdout
		setupOptions = origSetup
	}
}

func TestSetup_InvalidInput(t *testing.T) {
	var called []string
	options := []string{"Homebrew", "zsh", "git", "GitHub CLI", "asdf"}
	mockActions := make([]func(), len(options))
	for i, name := range options {
		name := name
		mockActions[i] = func() { called = append(called, name) }
	}

	input := bytes.NewBufferString("9,2,foo\n")
	output := &bytes.Buffer{}
	oldStdin := stdin
	oldStdout := stdout
	stdin = input
	stdout = output
	defer func() { stdin = oldStdin; stdout = oldStdout }()

	origSetup := setupOptions
	setupOptions = func() []struct{ name string; action func() } {
		opts := make([]struct{ name string; action func() }, len(options))
		for i := range options {
			opts[i] = struct{ name string; action func() }{options[i], mockActions[i]}
		}
		return opts
	}
	defer func() { setupOptions = origSetup }()

	Setup()

	if len(called) != 1 || called[0] != "zsh" {
		t.Errorf("Expected only zsh to be called, got %v", called)
	}
	if !strings.Contains(output.String(), "Invalid choice") {
		t.Errorf("Expected error message for invalid input")
	}
}

