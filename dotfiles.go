package main

import (
	"bufio"
	"fmt"
	"os"
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
This is a test
And this is the second line
Starting third line
$ `

func main() {
	printText(text)
	result := waitForInput()
	printText("Selected: " + result)
}
