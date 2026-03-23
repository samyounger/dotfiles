#!/usr/bin/env bash

set -euo pipefail

DOTFILES_REPO_SSH="${DOTFILES_REPO_SSH:-git@github.com:samyounger/dotfiles.git}"
DOTFILES_DIR="${DOTFILES_DIR:-$HOME/.dotfiles}"
ONEPASSWORD_SSH_AGENT_SOCK="${ONEPASSWORD_SSH_AGENT_SOCK:-$HOME/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock}"

log() {
  printf '\n==> %s\n' "$1"
}

fail() {
  printf 'Error: %s\n' "$1" >&2
  exit 1
}

ensure_macos() {
  if [ "$(uname -s)" != "Darwin" ]; then
    fail "This bootstrap script only supports macOS."
  fi
}

brew_bin() {
  if command -v brew >/dev/null 2>&1; then
    command -v brew
    return
  fi

  if [ -x /opt/homebrew/bin/brew ]; then
    printf '/opt/homebrew/bin/brew\n'
    return
  fi

  if [ -x /usr/local/bin/brew ]; then
    printf '/usr/local/bin/brew\n'
    return
  fi

  fail "Homebrew was expected but could not be found."
}

ensure_homebrew() {
  if command -v brew >/dev/null 2>&1 || [ -x /opt/homebrew/bin/brew ] || [ -x /usr/local/bin/brew ]; then
    log "Homebrew already installed"
    return
  fi

  log "Installing Homebrew"
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
}

install_base_packages() {
  local brew
  brew="$(brew_bin)"

  log "Installing git, zsh, asdf, and 1Password"
  "$brew" install git zsh asdf
  "$brew" install --cask 1password
}

ensure_1password_agent() {
  export SSH_AUTH_SOCK="$ONEPASSWORD_SSH_AGENT_SOCK"

  if [ -S "$SSH_AUTH_SOCK" ]; then
    log "1Password SSH agent already available"
    return
  fi

  log "Opening 1Password"
  open -a "1Password" || true

  cat <<EOF

1Password needs one manual step before the repo can be cloned with SSH:
  1. Sign in to 1Password.
  2. Unlock the app.
  3. Enable the SSH agent in 1Password settings.
  4. Make sure your GitHub SSH key is available there.

EOF

  while [ ! -S "$SSH_AUTH_SOCK" ]; do
    read -r -p "Press Enter after the 1Password SSH agent is ready..."
  done
}

verify_github_ssh() {
  local output

  log "Verifying GitHub SSH access"
  output="$(ssh -o StrictHostKeyChecking=accept-new -T git@github.com 2>&1 || true)"
  printf '%s\n' "$output"

  if ! printf '%s' "$output" | grep -q "successfully authenticated"; then
    fail "GitHub SSH authentication did not succeed."
  fi
}

clone_dotfiles_repo() {
  if [ -d "$DOTFILES_DIR/.git" ]; then
    log "Dotfiles repo already cloned at $DOTFILES_DIR"
    return
  fi

  if [ -e "$DOTFILES_DIR" ]; then
    fail "$DOTFILES_DIR exists but is not a git repository."
  fi

  log "Cloning dotfiles repo"
  git clone "$DOTFILES_REPO_SSH" "$DOTFILES_DIR"
}

copy_tool_versions() {
  local source_file target_file

  source_file="$DOTFILES_DIR/runcom/.tool-versions"
  target_file="$HOME/.tool-versions"

  if [ ! -f "$source_file" ]; then
    fail "Missing $source_file"
  fi

  if [ -e "$target_file" ]; then
    log "Keeping existing $target_file"
    return
  fi

  log "Copying .tool-versions into $HOME"
  cp "$source_file" "$target_file"
}

ensure_golang_plugin() {
  if asdf plugin list | grep -qx "golang"; then
    return
  fi

  log "Adding asdf golang plugin"
  asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
}

install_go_with_asdf() {
  local go_version

  go_version="$(awk '$1 == "golang" { print $2 }' "$DOTFILES_DIR/runcom/.tool-versions")"
  if [ -z "$go_version" ]; then
    fail "Could not determine the golang version from runcom/.tool-versions."
  fi

  ensure_golang_plugin

  log "Installing Go $go_version with asdf"
  asdf install golang "$go_version"
  asdf global golang "$go_version"
}

install_dotfiles() {
  log "Installing home dotfiles"
  cd "$DOTFILES_DIR"
  printf '2\n' | asdf exec go run .
}

print_next_steps() {
  cat <<EOF

Bootstrap complete.

Next steps:
  1. Open a new terminal window or run: exec zsh -l
  2. Re-enter the repo: cd "$DOTFILES_DIR"
  3. Run the CLI when you want the next setup stage:
     asdf exec go run .

Useful menu choices:
  1) Setup dependencies
  2) Install dotfiles
  3) Setup OSX Defaults

EOF
}

main() {
  ensure_macos
  ensure_homebrew
  install_base_packages
  ensure_1password_agent
  verify_github_ssh
  clone_dotfiles_repo
  copy_tool_versions
  install_go_with_asdf
  install_dotfiles
  print_next_steps
}

main "$@"
