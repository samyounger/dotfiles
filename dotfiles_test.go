package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestMain_Setup(t *testing.T) {
	var calledSetup bool
	setupFunc = func(input io.Reader, output io.Writer) { calledSetup = true }
	installFunc = func(output io.Writer) { t.Errorf("Should not call installFunc") }
	setupOSXDefaultsFunc = func(output io.Writer) { t.Errorf("Should not call setupOSXDefaultsFunc") }
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

func TestMain_InstallDotfiles(t *testing.T) {
	setupFunc = func(input io.Reader, output io.Writer) { t.Errorf("Should not call setupFunc") }
	var calledInstall bool
	installFunc = func(output io.Writer) { calledInstall = true }
	setupOSXDefaultsFunc = func(output io.Writer) { t.Errorf("Should not call setupOSXDefaultsFunc") }
	inputReader = strings.NewReader("2\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !calledInstall {
		t.Errorf("Expected installFunc to be called")
	}
}

func TestMain_SetupOSXDefaults(t *testing.T) {
	setupFunc = func(input io.Reader, output io.Writer) { t.Errorf("Should not call setupFunc") }
	installFunc = func(output io.Writer) { t.Errorf("Should not call installFunc") }
	var calledSetupOSXDefaults bool
	setupOSXDefaultsFunc = func(output io.Writer) { calledSetupOSXDefaults = true }
	inputReader = strings.NewReader("3\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !calledSetupOSXDefaults {
		t.Errorf("Expected setupOSXDefaultsFunc to be called")
	}
}

func TestMain_Exit(t *testing.T) {
	setupFunc = func(input io.Reader, output io.Writer) { t.Errorf("Should not call setupFunc") }
	installFunc = func(output io.Writer) { t.Errorf("Should not call installFunc") }
	setupOSXDefaultsFunc = func(output io.Writer) { t.Errorf("Should not call setupOSXDefaultsFunc") }
	inputReader = strings.NewReader("0\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !strings.Contains(buf.String(), "Exiting.") {
		t.Errorf("Expected 'Exiting.' to be printed")
	}
}

func TestMain_Default(t *testing.T) {
	setupFunc = func(input io.Reader, output io.Writer) { t.Errorf("Should not call setupFunc") }
	installFunc = func(output io.Writer) { t.Errorf("Should not call installFunc") }
	setupOSXDefaultsFunc = func(output io.Writer) { t.Errorf("Should not call setupOSXDefaultsFunc") }
	inputReader = strings.NewReader("foobar\n")
	var buf bytes.Buffer
	outputWriter = &buf

	main()

	if !strings.Contains(buf.String(), "Selected: foobar") {
		t.Errorf("Expected default handler to print input")
	}
}
