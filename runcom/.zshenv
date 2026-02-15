# ~/.zshenv example template
# Environment variables and settings for all zsh sessions

# Set user-specific environment variables
export EDITOR="vim"
export LANG="en_US.UTF-8"

# Add custom bin directories to PATH
export PATH="$HOME/.local/bin:$PATH"

# Add any other environment variables below
# export JAVA_HOME="$HOME/.sdkman/candidates/java/current"

export PIHOLE_SSH_USER=""
export PIHOLE_SSH_IP=""

# For Github-CLI authentication via 1Password
export SSH_AUTH_SOCK="$HOME/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"

export COURSIER_PATH="$HOME/Library/Application Support/Coursier/bin"
export GOPATH="$HOME/source/go"
export PYTHON_PATH="/usr/bin/python3"
