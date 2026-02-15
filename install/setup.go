package install

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/syounger-ut/dotfiles/install/dependencies"
)

// Setup interactively installs Homebrew, zsh, git, GitHub CLI, and asdf
func Dependencies() {
	options := []struct {
		name string
		action func()
	}{
		{"Homebrew", dependencies.InstallHomebrew},
		{"zsh", dependencies.InstallZsh},
		{"git", dependencies.InstallGit},
		{"GitHub CLI", dependencies.InstallGithubCLI},
		{"asdf", dependencies.InstallAsdf},
	}

	fmt.Println("Select which dependencies to install:")
	fmt.Println("0) All")
	for i, opt := range options {
		fmt.Printf("%d) %s\n", i+1, opt.name)
	}
	fmt.Print("Enter your choice (comma-separated for multiple, e.g. 1,3): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.TrimSpace(input)

	if input == "0" {
		for _, opt := range options {
			opt.action()
		}
		return
	}

	choices := strings.Split(input, ",")
	chosen := map[int]bool{}
	for _, c := range choices {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		idx, err := parseChoice(c, len(options))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid choice: %s\n", c)
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