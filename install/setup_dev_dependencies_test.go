package install

import (
	"bufio"
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

	DependenciesWithMocks(input, output, mockActions)

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

		DependenciesWithMocks(input, output, mockActions)

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

	DependenciesWithMocks(input, output, mockActions)

	if len(called) != 1 || called[0] != "zsh" {
		t.Errorf("Expected only zsh to be called, got %v", called)
	}
	if !strings.Contains(output.String(), "Invalid choice") {
		t.Errorf("Expected error message for invalid input")
	}
}

// DependenciesWithMocks is a test helper to inject mock actions for unit tests
func DependenciesWithMocks(input *bytes.Buffer, output *bytes.Buffer, actions []func()) {
	options := []struct {
		name   string
		action func()
	}{
		{"Homebrew", actions[0]},
		{"zsh", actions[1]},
		{"git", actions[2]},
		{"GitHub CLI", actions[3]},
		{"asdf", actions[4]},
	}

	output.WriteString("Select which dependencies to install:\n")
	output.WriteString("0) All\n")
	for _, opt := range options {
		output.WriteString(strings.TrimSpace(opt.name) + "\n")
	}
	output.WriteString("Enter your choice (comma-separated for multiple, e.g. 1,3): ")

	scanner := bufio.NewScanner(input)
	scanner.Scan()
	userInput := scanner.Text()
	userInput = strings.TrimSpace(userInput)

	if userInput == "0" {
		for _, opt := range options {
			opt.action()
		}
		return
	}

	choices := strings.Split(userInput, ",")
	chosen := map[int]bool{}
	for _, c := range choices {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		idx, err := parseChoice(c, len(options))
		if err != nil {
			output.WriteString("Invalid choice: " + c + "\n")
			continue
		}
		chosen[idx] = true
	}
	// Ensure deterministic order for tests
	var ordered []int
	for _, c := range choices {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		idx, err := parseChoice(c, len(options))
		if err == nil && chosen[idx] {
			ordered = append(ordered, idx)
			delete(chosen, idx)
		}
	}
	for _, idx := range ordered {
		options[idx].action()
	}
}
