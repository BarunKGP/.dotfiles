# ~/.zshrc: Main zsh shell configuration loader
# This file sources modular shell configuration from $DOTFILES/shell/

# Auto-detect DOTFILES location from symlink
# Handles both standard ~/.zshrc and non-standard locations
if [ -z "$DOTFILES" ]; then
    if [ -L "$HOME/.zshrc" ]; then
        export DOTFILES="$(dirname "$(dirname "$(readlink -f "$HOME/.zshrc")")")"
    fi
fi
export DOTFILES="${DOTFILES:-$HOME/.dotfiles}"

# Helper to source a file if it exists
_src() { [ -f "$1" ] && source "$1"; }

# Load modules in order
_src "$DOTFILES/shell/zsh/options.zsh"        # Zsh-specific options
_src "$DOTFILES/shell/common/env.sh"          # Environment detection (sets OS_TYPE)
_src "$DOTFILES/shell/common/path.sh"         # PATH management
_src "$DOTFILES/shell/common/aliases.sh"      # Common aliases
_src "$DOTFILES/shell/common/functions.sh"    # Common functions
_src "$DOTFILES/shell/common/os/$OS_TYPE.sh"  # OS-specific configuration
_src "$DOTFILES/shell/zsh/completion.zsh"     # Zsh completion
_src "$DOTFILES/shell/common/tools.sh"        # Tool initializations (nvm, pyenv, zoxide, etc.)
_src "$HOME/.config/shell/local.sh"           # Machine-local overrides (not tracked in git)

unset -f _src
