package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultComputerName    = "Neo"
	defaultLocale          = "en_GB@currency=GBP"
	defaultMeasurementUnit = "Centimeters"
)

var languagePrefs = []string{"en", "GB"}

func SetupOSXDefaults(output io.Writer) {
	if !runOSXDefaults(output) {
		return
	}
	fmt.Fprintln(output, "OSX defaults setup complete.")
}

func runOSXDefaults(output io.Writer) bool {
	if err := applyOSXDefaults(output); err != nil {
		fmt.Fprintf(output, "Failed to apply OSX defaults: %v\n", err)
		return false
	}
	return true
}

func applyOSXDefaults(output io.Writer) error {
	home, err := userHomeDir()
	if err != nil {
		return err
	}

	computerName := os.Getenv("COMPUTER_NAME")
	if strings.TrimSpace(computerName) == "" {
		computerName = defaultComputerName
	}

	screenshotFolder := filepath.Join(home, "Downloads")
	if err := os.MkdirAll(screenshotFolder, 0o755); err != nil {
		return fmt.Errorf("create screenshots folder: %w", err)
	}

	commands := buildOSXDefaultsCommands(computerName, screenshotFolder)
	for _, command := range commands {
		if err := runCommand(command.name, command.args...); err != nil {
			if command.allowFailure {
				fmt.Fprintf(output, "Ignoring failure for command: %s %s (%v)\n", command.name, strings.Join(command.args, " "), err)
				continue
			}
			return fmt.Errorf("%s %s: %w", command.name, strings.Join(command.args, " "), err)
		}
	}

	// Restart impacted apps so the new defaults take effect immediately.
	for _, app := range []string{"Address Book", "Calendar", "Contacts", "Dock", "Finder", "Mail", "Safari", "SystemUIServer", "iCal"} {
		if err := runCommand("killall", app); err != nil {
			fmt.Fprintf(output, "Ignoring failure for command: killall %s (%v)\n", app, err)
		}
	}

	return nil
}

type osxDefaultsCommand struct {
	name         string
	args         []string
	allowFailure bool
}

func buildOSXDefaultsCommands(computerName, screenshotFolder string) []osxDefaultsCommand {
	commands := []osxDefaultsCommand{
		// Close System Settings before changing preferences.
		{name: "osascript", args: []string{"-e", "tell application \"System Preferences\" to quit"}, allowFailure: true},
		// Prompt once for sudo so elevated commands can run.
		{name: "sudo", args: []string{"-v"}},

		// Computer & Host name.
		{name: "sudo", args: []string{"scutil", "--set", "ComputerName", computerName}},
		{name: "sudo", args: []string{"scutil", "--set", "HostName", computerName}},
		{name: "sudo", args: []string{"scutil", "--set", "LocalHostName", computerName}},
		{name: "sudo", args: []string{"defaults", "write", "/Library/Preferences/SystemConfiguration/com.apple.smb.server", "NetBIOSName", "-string", computerName}},

		// Localization.
		{name: "defaults", args: append([]string{"write", "NSGlobalDomain", "AppleLanguages", "-array"}, languagePrefs...)},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "AppleLocale", "-string", defaultLocale}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "AppleMeasurementUnits", "-string", defaultMeasurementUnit}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "AppleMetricUnits", "-bool", "true"}},
		{name: "sudo", args: []string{"defaults", "write", "/Library/Preferences/com.apple.timezone.auto", "Active", "-bool", "YES"}},
		// systemsetup can return Error:-99 on some machines while still applying.
		{name: "sudo", args: []string{"systemsetup", "-setusingnetworktime", "on"}, allowFailure: true},

		// System.
		{name: "sudo", args: []string{"systemsetup", "-setrestartfreeze", "on"}, allowFailure: true},
		{name: "sudo", args: []string{"pmset", "-a", "standbydelay", "86400"}},
		{name: "sudo", args: []string{"pmset", "-a", "sms", "0"}},
		{name: "defaults", args: []string{"write", "com.apple.sound.beep.feedback", "-bool", "false"}},
		{name: "sudo", args: []string{"nvram", "SystemAudioVolume= "}},
		{name: "sudo", args: []string{"nvram", "StartupMute=%01"}},
		{name: "defaults", args: []string{"write", "com.apple.menuextra.battery", "ShowPercent", "YES"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSAutomaticWindowAnimationsEnabled", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSWindowResizeTime", "-float", "0.001"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSNavPanelExpandedStateForSaveMode", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSNavPanelExpandedStateForSaveMode2", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "PMPrintingExpandedStateForPrint", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "PMPrintingExpandedStateForPrint2", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSDocumentSaveNewDocumentsToCloud", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.print.PrintingPrefs", "Quit When Finished", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.LaunchServices", "LSQuarantine", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.systempreferences", "NSQuitAlwaysKeepsWindows", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.CrashReporter", "DialogType", "-string", "none"}},
		{name: "launchctl", args: []string{"unload", "-w", "/System/Library/LaunchAgents/com.apple.notificationcenterui.plist"}, allowFailure: true},

		// Keyboard & Input.
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSAutomaticQuoteSubstitutionEnabled", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSAutomaticDashSubstitutionEnabled", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "AppleKeyboardUIMode", "-int", "3"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "ApplePressAndHoldEnabled", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "KeyRepeat", "-int", "1"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "InitialKeyRepeat", "-int", "15"}},
		{name: "defaults", args: []string{"write", "com.apple.BezelServices", "kDim", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.BezelServices", "kDimTime", "-int", "300"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "NSAutomaticSpellingCorrectionEnabled", "-bool", "false"}},

		// Trackpad, mouse, Bluetooth accessories.
		{name: "defaults", args: []string{"write", "com.apple.AppleMultitouchTrackpad", "Clicking", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "Clicking", "-bool", "true"}},
		{name: "defaults", args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.mouse.tapBehavior", "-int", "1"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "com.apple.mouse.tapBehavior", "-int", "1"}},
		{name: "defaults", args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "TrackpadCornerSecondaryClick", "-int", "2"}},
		{name: "defaults", args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "TrackpadRightClick", "-bool", "true"}},
		{name: "defaults", args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.trackpad.trackpadCornerClickBehavior", "-int", "1"}},
		{name: "defaults", args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.trackpad.enableSecondaryClick", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "AppleEnableSwipeNavigateWithScrolls", "-bool", "true"}},
		{name: "defaults", args: []string{"-currentHost", "write", "NSGlobalDomain", "com.apple.trackpad.threeFingerHorizSwipeGesture", "-int", "1"}},
		{name: "defaults", args: []string{"write", "com.apple.driver.AppleBluetoothMultitouch.trackpad", "TrackpadThreeFingerHorizSwipeGesture", "-int", "1"}},
		{name: "defaults", args: []string{"write", "com.apple.BluetoothAudioAgent", "Apple Bitpool Min (editable)", "-int", "40"}},

		// Screen.
		{name: "defaults", args: []string{"write", "com.apple.screensaver", "askForPassword", "-int", "1"}},
		{name: "defaults", args: []string{"write", "com.apple.screensaver", "askForPasswordDelay", "-int", "0"}},
		{name: "defaults", args: []string{"write", "com.apple.screencapture", "location", "-string", screenshotFolder}},
		{name: "defaults", args: []string{"write", "com.apple.screencapture", "type", "-string", "png"}},
		{name: "defaults", args: []string{"write", "com.apple.screencapture", "disable-shadow", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "AppleFontSmoothing", "-int", "2"}},

		// Finder.
		{name: "defaults", args: []string{"write", "com.apple.finder", "QuitMenuItem", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "DisableAllAnimations", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "AppleShowAllFiles", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "NSGlobalDomain", "AppleShowAllExtensions", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "ShowStatusBar", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "ShowPathbar", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "QLEnableTextSelection", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "_FXShowPosixPathInTitle", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "_FXSortFoldersFirst", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "FXDefaultSearchScope", "-string", "SCcf"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "FXEnableExtensionChangeWarning", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.desktopservices", "DSDontWriteNetworkStores", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.desktopservices", "DSDontWriteUSBStores", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.frameworks.diskimages", "skip-verify", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.frameworks.diskimages", "skip-verify-locked", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.frameworks.diskimages", "skip-verify-remote", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.NetworkBrowser", "BrowseAllInterfaces", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "FXPreferredViewStyle", "-string", "Nlsv"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "WarnOnEmptyTrash", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.finder", "FXInfoPanesExpanded", "-dict", "General", "-bool", "true", "OpenWith", "-bool", "true", "Privileges", "-bool", "true"}},

		// Dock.
		{name: "defaults", args: []string{"write", "com.apple.dock", "show-process-indicators", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "launchanim", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "autohide", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "showhidden", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "no-bouncing", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "wvous-tl-corner", "-int", "0"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "wvous-tr-corner", "-int", "0"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "wvous-bl-corner", "-int", "0"}},
		{name: "defaults", args: []string{"write", "com.apple.dock", "wvous-br-corner", "-int", "0"}},
		{name: "defaults", args: []string{"write", "com.Apple.Dock", "show-recents", "-bool", "false"}},

		// Mail.
		{name: "defaults", args: []string{"write", "com.apple.mail", "DraftsViewerAttributes", "-dict-add", "DisplayInThreadedMode", "-string", "yes"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "DisableReplyAnimations", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "DisableSendAnimations", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "AddressesIncludeNameOnPasteboard", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "DisableInlineAttachmentViewing", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "SpellCheckingBehavior", "-string", "NoSpellCheckingEnabled"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "MailSound", "-string", ""}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "PlayMailSounds", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "ConversationViewMarkAllAsRead", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "IndexTrash", "-bool", "false"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "AutoFetch", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "PollTime", "-string", "-1"}},
		{name: "defaults", args: []string{"write", "com.apple.mail", "ConversationViewSortDescending", "-bool", "true"}},

		// Calendar.
		{name: "defaults", args: []string{"write", "com.apple.iCal", "Show Week Numbers", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.iCal", "first day of week", "-int", "1"}},

		// Terminal.
		{name: "defaults", args: []string{"write", "com.apple.terminal", "StringEncodings", "-array", "4"}},
		{name: "defaults", args: []string{"write", "com.apple.terminal", "Default Window Settings", "-string", "Pro"}},
		{name: "defaults", args: []string{"write", "com.apple.terminal", "Startup Window Settings", "-string", "Pro"}},
		{name: "defaults", args: []string{"write", "com.apple.Terminal", "ShowLineMarks", "-int", "0"}},

		// Activity Monitor.
		{name: "defaults", args: []string{"write", "com.apple.ActivityMonitor", "OpenMainWindow", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.ActivityMonitor", "IconType", "-int", "5"}},
		{name: "defaults", args: []string{"write", "com.apple.ActivityMonitor", "ShowCategory", "-int", "0"}},
		{name: "defaults", args: []string{"write", "com.apple.ActivityMonitor", "SortColumn", "-string", "CPUUsage"}},
		{name: "defaults", args: []string{"write", "com.apple.ActivityMonitor", "SortDirection", "-int", "0"}},

		// Software updates.
		{name: "defaults", args: []string{"write", "com.apple.SoftwareUpdate", "AutomaticCheckEnabled", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.SoftwareUpdate", "ScheduleFrequency", "-string", "7"}},
		{name: "defaults", args: []string{"write", "com.apple.SoftwareUpdate", "AutomaticDownload", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.SoftwareUpdate", "CriticalUpdateInstall", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.commerce", "AutoUpdate", "-bool", "true"}},
		{name: "defaults", args: []string{"write", "com.apple.commerce", "AutoUpdateRestartRequired", "-bool", "true"}},
	}

	return commands
}
