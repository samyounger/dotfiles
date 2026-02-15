package application_settings

func BuildChromeDefaultsCommands() []SettingsCommand {
	return []SettingsCommand{
		// Disable swipe navigation.
		{Name: "defaults", Args: []string{"write", "com.google.Chrome", "AppleEnableSwipeNavigateWithScrolls", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.google.Chrome.canary", "AppleEnableSwipeNavigateWithScrolls", "-bool", "false"}},

		// Use the system-native print preview dialog.
		{Name: "defaults", Args: []string{"write", "com.google.Chrome", "DisablePrintPreview", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.google.Chrome.canary", "DisablePrintPreview", "-bool", "true"}},

		// Expand the print dialog by default.
		{Name: "defaults", Args: []string{"write", "com.google.Chrome", "PMPrintingExpandedStateForPrint2", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.google.Chrome.canary", "PMPrintingExpandedStateForPrint2", "-bool", "true"}},
	}
}
