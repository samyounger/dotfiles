# ~/.zshenv example template
# Environment variables and settings for all zsh sessions

# Set user-specific environment variables
export EDITOR="nvim"

## Set language environment
export LANG="en_GB.UTF-8"

# Add custom bin directories to PATH
export PATH="$HOME/.local/bin:$PATH"

# For Github-CLI authentication via 1Password
export SSH_AUTH_SOCK="$HOME/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"

# Path to Github MCP
export COPILOT_MCP_CONFIG="$HOME/.config/copilot/mcp.json"

# Path to aws cli config
export AWS_CONFIG_FILE="$HOME/.config/aws/config"
export AWS_PROFILE="engineering"

# GPG key
export GPG_TTY=$(tty)
