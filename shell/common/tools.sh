# Initialize development tools (nvm, pyenv, zoxide, etc.)

# NVM - Node Version Manager
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"
if [ -n "$BASH_VERSION" ] && [ -s "$NVM_DIR/bash_completion" ]; then
    . "$NVM_DIR/bash_completion"
fi

# pyenv - Python Version Manager
export PYENV_ROOT="$HOME/.pyenv"
[ -d "$PYENV_ROOT/bin" ] && export PATH="$PYENV_ROOT/bin:$PATH"
if command -v pyenv >/dev/null 2>&1; then
    eval "$(pyenv init -)"
fi

# zoxide - Smarter cd replacement
if command -v zoxide >/dev/null 2>&1; then
    if [ -n "$BASH_VERSION" ]; then
        eval "$(zoxide init bash)"
    elif [ -n "$ZSH_VERSION" ]; then
        eval "$(zoxide init zsh)"
    fi
    alias cd='z'
    alias cdi='zi'
fi

# .local/bin environment
[ -f "$HOME/.local/bin/env" ] && . "$HOME/.local/bin/env"
