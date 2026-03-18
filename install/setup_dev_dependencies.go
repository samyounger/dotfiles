package install

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/samyounger/dotfiles/install/dependencies"
)

// Dependencies interactively installs Homebrew, zsh, git, GitHub CLI, and asdf
func Dependencies(input io.Reader, output io.Writer) {
	options := []struct {
		name   string
		action func()
	}{
		{"Homebrew", dependencies.InstallHomebrew},
		{"zsh", dependencies.InstallZsh},
		{"git", dependencies.InstallGit},
		{"GitHub CLI", dependencies.InstallGithubCLI},
		{"asdf", dependencies.InstallAsdf},
	}

	fmt.Fprintln(output, "Select which dependencies to install:")
	fmt.Fprintln(output, "0) All")
	for i, opt := range options {
		fmt.Fprintf(output, "%d) %s\n", i+1, opt.name)
	}
	fmt.Fprint(output, "Enter your choice (comma-separated for multiple, e.g. 1,3): ")

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
			fmt.Fprintf(output, "Invalid choice: %s\n", c)
			continue
		}
		chosen[idx] = true
	}
	for idx := range chosen {
		options[idx].action()
	}
}

func parseChoice(s string, max int) (int, error) {
	switch s {
	case "1", "2", "3", "4", "5":
		idx := int(s[0] - '1')
		if idx >= 0 && idx < max {
			return idx, nil
		}
	}
	return 0, fmt.Errorf("out of range")
}
