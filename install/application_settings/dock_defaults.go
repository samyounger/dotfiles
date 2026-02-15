package application_settings

var desiredDockApplications = []string{
	"/System/Applications/Apps.app",
	"/Applications/Google Chrome.app",
	"/Applications/Slack.app",
	"/Applications/iTerm.app",
	"/Applications/ChatGPT.app",
	"/System/Applications/Music.app",
}

func DesiredDockApplications() []string {
	apps := make([]string, len(desiredDockApplications))
	copy(apps, desiredDockApplications)
	return apps
}

func BuildDockFindCommand(appPath string) SettingsCommand {
	return SettingsCommand{Name: "dockutil", Args: []string{"--find", appPath}}
}

func BuildDockAddCommand(appPath string) SettingsCommand {
	return SettingsCommand{Name: "dockutil", Args: []string{"--no-restart", "--add", appPath}}
}

func BuildDockRestartCommand() SettingsCommand {
	return SettingsCommand{Name: "killall", Args: []string{"Dock"}}
}
