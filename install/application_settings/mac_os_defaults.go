package application_settings

const (
	defaultLocale          = "en_GB@currency=GBP"
	defaultMeasurementUnit = "Centimeters"
)

var languagePrefs = []string{"en", "GB"}

type MacOSDefaultsCommand struct {
	Name         string
	Args         []string
	AllowFailure bool
}

func BuildMacOSDefaultsCommands(computerName, screenshotFolder string) []MacOSDefaultsCommand {
	commands := []MacOSDefaultsCommand{
		// Close System Settings before changing preferences.
		{Name: "osascript", Args: []string{"-e", "tell application \"System Preferences\" to quit"}, AllowFailure: true},
		// Prompt once for sudo so elevated commands can run.
		{Name: "sudo", Args: []string{"-v"}},

		// Computer & Host name.
		{Name: "sudo", Args: []string{"scutil", "--set", "ComputerName", computerName}},
		{Name: "sudo", Args: []string{"scutil", "--set", "HostName", computerName}},
		{Name: "sudo", Args: []string{"scutil", "--set", "LocalHostName", computerName}},
		{Name: "sudo", Args: []string{"defaults", "write", "/Library/Preferences/SystemConfiguration/com.apple.smb.server", "NetBIOSName", "-string", computerName}},

		// Localization.
		{Name: "defaults", Args: append([]string{"write", "NSGlobalDomain", "AppleLanguages", "-array"}, languagePrefs...)},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "AppleLocale", "-string", defaultLocale}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "AppleMeasurementUnits", "-string", defaultMeasurementUnit}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "AppleMetricUnits", "-bool", "true"}},
		{Name: "sudo", Args: []string{"defaults", "write", "/Library/Preferences/com.apple.timezone.auto", "Active", "-bool", "YES"}},
		// systemsetup can return Error:-99 on some machines while still applying.
		{Name: "sudo", Args: []string{"systemsetup", "-setusingnetworktime", "on"}, AllowFailure: true},

		// System.
		{Name: "sudo", Args: []string{"systemsetup", "-setrestartfreeze", "on"}, AllowFailure: true},
		{Name: "sudo", Args: []string{"pmset", "-a", "standbydelay", "86400"}},
		{Name: "sudo", Args: []string{"pmset", "-a", "sms", "0"}},
		{Name: "defaults", Args: []string{"write", "com.apple.sound.beep.feedback", "-bool", "false"}},
		{Name: "sudo", Args: []string{"nvram", "SystemAudioVolume= "}},
		{Name: "sudo", Args: []string{"nvram", "StartupMute=%01"}},
		{Name: "defaults", Args: []string{"write", "com.apple.menuextra.battery", "ShowPercent", "YES"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSAutomaticWindowAnimationsEnabled", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSWindowResizeTime", "-float", "0.001"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSNavPanelExpandedStateForSaveMode", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSNavPanelExpandedStateForSaveMode2", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "PMPrintingExpandedStateForPrint", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "PMPrintingExpandedStateForPrint2", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSDocumentSaveNewDocumentsToCloud", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.print.PrintingPrefs", "Quit When Finished", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.LaunchServices", "LSQuarantine", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.systempreferences", "NSQuitAlwaysKeepsWindows", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.CrashReporter", "DialogType", "-string", "none"}},
		{Name: "launchctl", Args: []string{"unload", "-w", "/System/Library/LaunchAgents/com.apple.notificationcenterui.plist"}, AllowFailure: true},

		// Keyboard & Input.
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSAutomaticQuoteSubstitutionEnabled", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSAutomaticDashSubstitutionEnabled", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "AppleKeyboardUIMode", "-int", "3"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "ApplePressAndHoldEnabled", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "KeyRepeat", "-int", "1"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "InitialKeyRepeat", "-int", "15"}},
		{Name: "defaults", Args: []string{"write", "com.apple.BezelServices", "kDim", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.BezelServices", "kDimTime", "-int", "300"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "NSAutomaticSpellingCorrectionEnabled", "-bool", "false"}},

		// Trackpad, mouse, Bluetooth accessories.
		{Name: "defaults", Args: []string{"write", "com.apple.AppleMultitouchTrackpad", "Clicking", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "Clicking", "-bool", "true"}},
		{Name: "defaults", Args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.mouse.tapBehavior", "-int", "1"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "com.apple.mouse.tapBehavior", "-int", "1"}},
		{Name: "defaults", Args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "TrackpadCornerSecondaryClick", "-int", "2"}},
		{Name: "defaults", Args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "TrackpadRightClick", "-bool", "true"}},
		{Name: "defaults", Args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.trackpad.trackpadCornerClickBehavior", "-int", "1"}},
		{Name: "defaults", Args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.trackpad.enableSecondaryClick", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "AppleEnableSwipeNavigateWithScrolls", "-bool", "true"}},
		{Name: "defaults", Args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.trackpad.threeFingerHorizSwipeGesture", "-int", "1"}},
		{Name: "defaults", Args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "TrackpadThreeFingerHorizSwipeGesture", "-int", "1"}},
		{Name: "defaults", Args: []string{"write", "com.apple.BluetoothAudioAgent", "Apple Bitpool Min (editable)", "-int", "40"}},

		// Screen.
		{Name: "defaults", Args: []string{"write", "com.apple.screensaver", "askForPassword", "-int", "1"}},
		{Name: "defaults", Args: []string{"write", "com.apple.screensaver", "askForPasswordDelay", "-int", "0"}},
		{Name: "defaults", Args: []string{"write", "com.apple.screencapture", "location", "-string", screenshotFolder}},
		{Name: "defaults", Args: []string{"write", "com.apple.screencapture", "type", "-string", "png"}},
		{Name: "defaults", Args: []string{"write", "com.apple.screencapture", "disable-shadow", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "AppleFontSmoothing", "-int", "2"}},

		// Finder.
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "QuitMenuItem", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "DisableAllAnimations", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "AppleShowAllFiles", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "NSGlobalDomain", "AppleShowAllExtensions", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "ShowStatusBar", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "ShowPathbar", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "QLEnableTextSelection", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "_FXShowPosixPathInTitle", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "_FXSortFoldersFirst", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "FXDefaultSearchScope", "-string", "SCcf"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "FXEnableExtensionChangeWarning", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.desktopservices", "DSDontWriteNetworkStores", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.desktopservices", "DSDontWriteUSBStores", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.frameworks.diskimages", "skip-verify", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.frameworks.diskimages", "skip-verify-locked", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.frameworks.diskimages", "skip-verify-remote", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.NetworkBrowser", "BrowseAllInterfaces", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "FXPreferredViewStyle", "-string", "Nlsv"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "WarnOnEmptyTrash", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.finder", "FXInfoPanesExpanded", "-dict", "General", "-bool", "true", "OpenWith", "-bool", "true", "Privileges", "-bool", "true"}},

		// Dock.
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "show-process-indicators", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "launchanim", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "autohide", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "showhidden", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "no-bouncing", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "wvous-tl-corner", "-int", "0"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "wvous-tr-corner", "-int", "0"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "wvous-bl-corner", "-int", "0"}},
		{Name: "defaults", Args: []string{"write", "com.apple.dock", "wvous-br-corner", "-int", "0"}},
		{Name: "defaults", Args: []string{"write", "com.Apple.Dock", "show-recents", "-bool", "false"}},

		// Mail.
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "DraftsViewerAttributes", "-dict-add", "DisplayInThreadedMode", "-string", "yes"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "DisableReplyAnimations", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "DisableSendAnimations", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "AddressesIncludeNameOnPasteboard", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "DisableInlineAttachmentViewing", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "SpellCheckingBehavior", "-string", "NoSpellCheckingEnabled"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "MailSound", "-string", ""}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "PlayMailSounds", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "ConversationViewMarkAllAsRead", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "IndexTrash", "-bool", "false"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "AutoFetch", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "PollTime", "-string", "-1"}},
		{Name: "defaults", Args: []string{"write", "com.apple.mail", "ConversationViewSortDescending", "-bool", "true"}},

		// Calendar.
		{Name: "defaults", Args: []string{"write", "com.apple.iCal", "Show Week Numbers", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.iCal", "first day of week", "-int", "1"}},

		// Terminal.
		{Name: "defaults", Args: []string{"write", "com.apple.terminal", "StringEncodings", "-array", "4"}},
		{Name: "defaults", Args: []string{"write", "com.apple.terminal", "Default Window Settings", "-string", "Pro"}},
		{Name: "defaults", Args: []string{"write", "com.apple.terminal", "Startup Window Settings", "-string", "Pro"}},
		{Name: "defaults", Args: []string{"write", "com.apple.Terminal", "ShowLineMarks", "-int", "0"}},

		// Activity Monitor.
		{Name: "defaults", Args: []string{"write", "com.apple.ActivityMonitor", "OpenMainWindow", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.ActivityMonitor", "IconType", "-int", "5"}},
		{Name: "defaults", Args: []string{"write", "com.apple.ActivityMonitor", "ShowCategory", "-int", "0"}},
		{Name: "defaults", Args: []string{"write", "com.apple.ActivityMonitor", "SortColumn", "-string", "CPUUsage"}},
		{Name: "defaults", Args: []string{"write", "com.apple.ActivityMonitor", "SortDirection", "-int", "0"}},

		// Software updates.
		{Name: "defaults", Args: []string{"write", "com.apple.SoftwareUpdate", "AutomaticCheckEnabled", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.SoftwareUpdate", "ScheduleFrequency", "-string", "7"}},
		{Name: "defaults", Args: []string{"write", "com.apple.SoftwareUpdate", "AutomaticDownload", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.SoftwareUpdate", "CriticalUpdateInstall", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.commerce", "AutoUpdate", "-bool", "true"}},
		{Name: "defaults", Args: []string{"write", "com.apple.commerce", "AutoUpdateRestartRequired", "-bool", "true"}},
	}

	return commands
}
