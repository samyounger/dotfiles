package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"github.com/syounger-ut/dotfiles/install/dependencies"
)

// Dependency injection for testability
var installGitFunc = dependencies.InstallGit
var installAsdfFunc = dependencies.InstallAsdf
var inputReader io.Reader = os.Stdin
var outputWriter io.Writer = os.Stdout

func printText(text string) {
	fmt.Fprint(outputWriter, text)
}

func waitForInput() string {
	input := bufio.NewScanner(inputReader)
	input.Scan()
	return input.Text()
}

const text = `
Options:
1) Install git
2) Install asdf
3) Exit
$ `

func main() {
	printText(text)
	result := waitForInput()
	switch result {
	case "1":
		installGitFunc()
	case "2":
		installAsdfFunc()
	case "3":
		fmt.Fprintln(outputWriter, "Exiting.")
		return
	default:
		printText("Selected: " + result + "\n")
	}
}
