# dotfiles
My collection of dotfiles

## Blank-slate macOS bootstrap

On a relatively fresh Mac, run the bootstrap script first. It installs the minimum needed to clone and use this repo:

- Homebrew
- `git`
- `zsh`
- `asdf`
- `1Password`
- Go via `asdf`

```bash
curl -fsSL https://raw.githubusercontent.com/samyounger/dotfiles/main/install/bootstrap_macos.sh | bash
```

The bootstrap flow will pause for the manual `1Password` step so you can sign in, unlock it, and expose your GitHub SSH key through the `1Password` SSH agent before the repo is cloned.

## Setup

After bootstrap, the repository should already exist at `~/.dotfiles`, the core home files should be installed, and the Go CLI should be runnable.

For contributors to this repository, install `gitleaks` through the Go CLI dependency menu, then enable the repo-managed Git hooks:

```bash
cd "$HOME/.dotfiles"
asdf exec go run .
# choose 6) gitleaks, or 0) All
./install/setup_git_hooks.sh
```

This repository uses `gitleaks` for free local and CI secrets scanning. The local hooks run on both commit and push, and GitHub Actions runs the same category of scan on pull requests and pushes using the pinned `gitleaks` CLI in Docker. To make the GitHub Action a true merge blocker, mark `Secret Scan / gitleaks` as a required status check in branch protection.

To re-run the installer manually:

```bash
cd "$HOME/.dotfiles"
asdf exec go run .
```

The main flows are:

- `1` Setup dependencies
- `2` Install dotfiles
- `3` Setup macOS defaults

```
 $ asdf exec go run .
```

## Documentation

- `Architecture.md` describes the repository layout, package responsibilities, diagrams, and coding conventions.
- `Agent.md` documents the preferred contributor workflow, including tests-first development and final verification.

## References

The original implementation of this repository is based on https://github.com/webpro/dotfiles/blob/main, following this post https://www.webpro.nl/articles/getting-started-with-dotfiles.

Discovered while looking at the Github guide to dotfiles https://dotfiles.github.io/tutorials/.
