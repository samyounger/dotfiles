# ~/.zshenv example template
# Environment variables and settings for all zsh sessions

# Set user-specific environment variables
export EDITOR="nvim"
export LANG="en_GB.UTF-8"

# Add custom bin directories to PATH
export PATH="$HOME/.local/bin:$PATH"

# For Github-CLI authentication via 1Password
export SSH_AUTH_SOCK="$HOME/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"

# Source local overrides if the file exists
source "$HOME/.zshenv.local" 2>/dev/null