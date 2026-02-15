package install

import (
	"bytes"
	"strings"
	"testing"
)

func TestDependencies_All(t *testing.T) {
	var called []string
	options := []string{"Homebrew", "zsh", "git", "GitHub CLI", "asdf"}
	mockActions := make([]func(), len(options))
	for i, name := range options {
		name := name
		mockActions[i] = func() { called = append(called, name) }
	}

	input := bytes.NewBufferString("0\n")
	output := &bytes.Buffer{}

	Dependencies(input, output)

	if len(called) != len(options) {
		t.Errorf("Expected all actions to be called, got %v", called)
	}
}

func TestDependencies_SingleAndMultiple(t *testing.T) {
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

		Dependencies(input, output)

		if len(called) != len(tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, called)
		}
		for i, v := range tc.expected {
			if called[i] != v {
				t.Errorf("Expected %v, got %v", tc.expected, called)
			}
		}
	}
}

func TestDependencies_InvalidInput(t *testing.T) {
	var called []string
	options := []string{"Homebrew", "zsh", "git", "GitHub CLI", "asdf"}
	mockActions := make([]func(), len(options))
	for i, name := range options {
		name := name
		mockActions[i] = func() { called = append(called, name) }
	}

	input := bytes.NewBufferString("9,2,foo\n")
	output := &bytes.Buffer{}

	Dependencies(input, output)

	if len(called) != 1 || called[0] != "zsh" {
		t.Errorf("Expected only zsh to be called, got %v", called)
	}
	if !strings.Contains(output.String(), "Invalid choice") {
		t.Errorf("Expected error message for invalid input")
	}
}
