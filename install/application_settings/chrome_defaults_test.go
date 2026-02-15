package application_settings

import "testing"

func TestBuildChromeDefaultsCommands_ContainsAllExpectedCommands(t *testing.T) {
	commands := BuildChromeDefaultsCommands()
	if len(commands) != 6 {
		t.Fatalf("expected 6 chrome defaults commands, got %d", len(commands))
	}

	expected := []struct {
		name string
		args []string
	}{
		{"defaults", []string{"write", "com.google.Chrome", "AppleEnableSwipeNavigateWithScrolls", "-bool", "false"}},
		{"defaults", []string{"write", "com.google.Chrome.canary", "AppleEnableSwipeNavigateWithScrolls", "-bool", "false"}},
		{"defaults", []string{"write", "com.google.Chrome", "DisablePrintPreview", "-bool", "true"}},
		{"defaults", []string{"write", "com.google.Chrome.canary", "DisablePrintPreview", "-bool", "true"}},
		{"defaults", []string{"write", "com.google.Chrome", "PMPrintingExpandedStateForPrint2", "-bool", "true"}},
		{"defaults", []string{"write", "com.google.Chrome.canary", "PMPrintingExpandedStateForPrint2", "-bool", "true"}},
	}

	for _, want := range expected {
		if !hasMacOSDefaultCommand(commands, want.name, want.args...) {
			t.Fatalf("missing command: %s %v", want.name, want.args)
		}
	}
}

func TestBuildChromeDefaultsCommands_NoAllowFailureCommands(t *testing.T) {
	commands := BuildChromeDefaultsCommands()
	for _, command := range commands {
		if command.AllowFailure {
			t.Fatalf("did not expect allow-failure chrome command: %s %v", command.Name, command.Args)
		}
	}
}
