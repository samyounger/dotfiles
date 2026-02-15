package application_settings

import "testing"

func TestDesiredDockApplications_ContainsExpectedApplications(t *testing.T) {
	apps := DesiredDockApplications()
	if len(apps) != 6 {
		t.Fatalf("expected 6 dock applications, got %d", len(apps))
	}

	expected := []string{
		"/System/Applications/Apps.app",
		"/Applications/Google Chrome.app",
		"/Applications/Slack.app",
		"/Applications/iTerm.app",
		"/Applications/ChatGPT.app",
		"/System/Applications/Music.app",
	}
	for _, app := range expected {
		found := false
		for _, actual := range apps {
			if actual == app {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected application missing: %s", app)
		}
	}
}

func TestBuildDockFindCommand(t *testing.T) {
	cmd := BuildDockFindCommand("/Applications/Google Chrome.app")
	if cmd.Name != "dockutil" {
		t.Fatalf("expected dockutil command, got %s", cmd.Name)
	}
	if !hasSettingsCommand([]SettingsCommand{cmd}, "dockutil", "--find", "/Applications/Google Chrome.app") {
		t.Fatalf("unexpected find command args: %v", cmd.Args)
	}
}

func TestBuildDockAddCommand(t *testing.T) {
	cmd := BuildDockAddCommand("/Applications/Google Chrome.app")
	if cmd.Name != "dockutil" {
		t.Fatalf("expected dockutil command, got %s", cmd.Name)
	}
	if !hasSettingsCommand([]SettingsCommand{cmd}, "dockutil", "--no-restart", "--add", "/Applications/Google Chrome.app") {
		t.Fatalf("unexpected add command args: %v", cmd.Args)
	}
}

func TestBuildDockRestartCommand(t *testing.T) {
	cmd := BuildDockRestartCommand()
	if cmd.Name != "killall" {
		t.Fatalf("expected killall command, got %s", cmd.Name)
	}
	if !hasSettingsCommand([]SettingsCommand{cmd}, "killall", "Dock") {
		t.Fatalf("unexpected restart command args: %v", cmd.Args)
	}
}

func TestDockCommands_NoAllowFailure(t *testing.T) {
	commands := []SettingsCommand{
		BuildDockFindCommand("/Applications/Google Chrome.app"),
		BuildDockAddCommand("/Applications/Google Chrome.app"),
		BuildDockRestartCommand(),
	}
	for _, command := range commands {
		if command.AllowFailure {
			t.Fatalf("did not expect allow-failure dock command: %s %v", command.Name, command.Args)
		}
	}
}
