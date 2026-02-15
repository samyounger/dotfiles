package install

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestSetupOSXDefaults_Success(t *testing.T) {
	defer resetDotfilesMocks()

	runCommand = func(name string, args ...string) error { return nil }

	out := &bytes.Buffer{}
	SetupOSXDefaults(out)

	if !strings.Contains(out.String(), "OSX defaults setup complete.") {
		t.Fatalf("expected success message, got: %s", out.String())
	}
}

func TestSetupOSXDefaults_Failure(t *testing.T) {
	defer resetDotfilesMocks()

	runCommand = func(name string, args ...string) error { return errors.New("boom") }

	out := &bytes.Buffer{}
	SetupOSXDefaults(out)

	if !strings.Contains(out.String(), "Failed to apply OSX defaults") {
		t.Fatalf("expected osxdefaults error message, got: %s", out.String())
	}
}
