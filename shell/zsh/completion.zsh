# Zsh completion configuration

# Load completion system
autoload -Uz compinit && compinit

# Completion styling
zstyle ':completion:*' menu select
zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}'
