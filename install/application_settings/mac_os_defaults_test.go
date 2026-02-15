package application_settings

import "testing"

func TestBuildMacOSDefaultsCommands_ContainsCoreCommands(t *testing.T) {
	commands := BuildMacOSDefaultsCommands("NeoTest", "/tmp/screens")
	if len(commands) == 0 {
		t.Fatal("expected command list to be non-empty")
	}

	if !hasMacOSDefaultCommand(commands, "sudo", "scutil", "--set", "ComputerName", "NeoTest") {
		t.Fatalf("expected computer name command in command list")
	}

	if !hasMacOSDefaultCommand(commands, "defaults", "write", "com.apple.screencapture", "location", "-string", "/tmp/screens") {
		t.Fatalf("expected screenshot location command in command list")
	}
}

func TestBuildMacOSDefaultsCommands_AllowFailureMarkers(t *testing.T) {
	commands := BuildMacOSDefaultsCommands("NeoTest", "/tmp/screens")

	if !hasMacOSDefaultCommandAllowFailure(commands, "osascript", []string{"-e", "tell application \"System Preferences\" to quit"}) {
		t.Fatalf("expected osascript command to be allow-failure")
	}

	if !hasMacOSDefaultCommandAllowFailure(commands, "sudo", []string{"systemsetup", "-setusingnetworktime", "on"}) {
		t.Fatalf("expected network time command to be allow-failure")
	}
}

func hasMacOSDefaultCommand(commands []MacOSDefaultsCommand, name string, args ...string) bool {
	for _, command := range commands {
		if command.Name != name {
			continue
		}
		if len(command.Args) != len(args) {
			continue
		}
		matched := true
		for i := range args {
			if command.Args[i] != args[i] {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

func hasMacOSDefaultCommandAllowFailure(commands []MacOSDefaultsCommand, name string, args []string) bool {
	for _, command := range commands {
		if command.Name != name {
			continue
		}
		if len(command.Args) != len(args) {
			continue
		}
		matched := true
		for i := range args {
			if command.Args[i] != args[i] {
				matched = false
				break
			}
		}
		if matched && command.AllowFailure {
			return true
		}
	}
	return false
}
