# dotfiles
My collection of dotfiles

## Setup

Place this in the `.zshrc` file

```sh
for DOTFILE in `find "$HOME/development/05-misc/.dotfiles" -not -path '*.md' -not -path '*.git/*'`
do
  [ -f "$DOTFILE" ] && source "$DOTFILE"
done
```

## References

The original implementation of this repository is based on https://github.com/webpro/dotfiles/blob/main, following this post https://www.webpro.nl/articles/getting-started-with-dotfiles.

Discovered while looking at the Github guide to dotfiles https://dotfiles.github.io/tutorials/.
