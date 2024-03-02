## Setup default OSX Settings
source ./osxdefaults.sh

## Symlink to home files
ln -sv "$HOME/development/05-misc/runcom/.zshrc" ~
ln -sv "$HOME/.dotfiles/runcom/.zshenv" ~
ln -sv "$HOME~/.dotfiles/runcom/.inputrc" ~
ln -sv "$HOME/.dotfiles/git/.gitconfig" ~
ln -sv "$HOME/.dotfiles/git/.gitignore_global" ~

## ASDF installation
## https://asdf-vm.com/guide/getting-started.html
function installASDF () {
  ## Install dependencies for ASDF
  brew install coreutils curl git

  # Install ASDF
  git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.13.1
}

installASDF

echo "🎉 Successful installation 🙌. Lets get started 🚀"
