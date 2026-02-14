package main

import (
	"bufio"
	"fmt"
	"os"
	"github.com/syounger-ut/dotfiles/install/dependencies"
)

func printText(text string) {
	fmt.Print(text)
}

func waitForInput() string {
	input := bufio.NewScanner(os.Stdin)
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
		dependencies.InstallGit()
	case "2":
		dependencies.InstallAsdf()
	case "3":
		fmt.Println("Exiting.")
		return
	default:
		printText("Selected: " + result + "\n")
	}
}
