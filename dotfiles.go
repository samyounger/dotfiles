package main

import (
	"bufio"
	"fmt"
	"github.com/syounger-ut/dotfiles/install"
	"io"
	"os"
)

// Dependency injection for testability
var setupFunc = install.Dependencies
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
1) Setup (dependencies)
5) Exit
$ `

func main() {
	printText(text)
	result := waitForInput()
	switch result {
	case "1":
		setupFunc()
	case "5":
		fmt.Fprintln(outputWriter, "Exiting.")
		return
	default:
		printText("Selected: " + result + "\n")
	}
}
