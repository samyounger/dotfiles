package install

import (
	"github.com/syounger-ut/dotfiles/install/dependencies"
)

// Setup installs Homebrew and zsh
func Setup() {
	dependencies.InstallHomebrew()
	dependencies.InstallZsh()
	dependencies.InstallGit()
	dependencies.InstallGithubCLI()
	dependencies.InstallAsdf()
}