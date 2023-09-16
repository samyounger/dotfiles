#!/bin/sh

dockutil --no-restart --remove all
dockutil --no-restart --add "/Applications/Google Chrome.app"
dockutil --no-restart --add "/System/Applications/Slack.app"
dockutil --no-restart --add "/System/Applications/Utilities/iTerm.app"

killall Dock
