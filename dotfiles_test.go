package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestMain_Setup(t *testing.T) {
	var calledSetup bool
	setupFunc = func() { calledSetup = true }
	inputReader = strings.NewReader("1\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !calledSetup {
		t.Errorf("Expected setupFunc to be called")
	}
	if !strings.Contains(buf.String(), "Options:") {
		t.Errorf("Expected menu to be printed")
	}
}

func TestMain_Exit(t *testing.T) {
	setupFunc = func() { t.Errorf("Should not call setupFunc") }
	inputReader = strings.NewReader("5\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !strings.Contains(buf.String(), "Exiting.") {
		t.Errorf("Expected 'Exiting.' to be printed")
	}
}

func TestMain_Default(t *testing.T) {
	setupFunc = func() { t.Errorf("Should not call setupFunc") }
	inputReader = strings.NewReader("foobar\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !strings.Contains(buf.String(), "Selected: foobar") {
		t.Errorf("Expected default handler to print input")
	}
}
