package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/samyounger/dotfiles/install"
)

// Dependency injection for testability
var setupFunc = install.Dependencies
var installFunc = install.InstallDotfiles
var setupMachineFunc = install.SetupMachine
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
1) Setup dependencies
2) Install dotfiles
3) Setup OSX Defaults
0) Exit
$ `

func main() {
	printText(text)
	result := waitForInput()
	switch result {
	case "1":
		setupFunc(inputReader, outputWriter)
	case "2":
		installFunc(outputWriter)
	case "3":
		setupMachineFunc(inputReader, outputWriter)
	case "0":
		fmt.Fprintln(outputWriter, "Exiting.")
		return
	default:
		printText("Selected: " + result + "\n")
	}
}
