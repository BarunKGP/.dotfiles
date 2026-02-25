# Cross-platform shell aliases

# Color support for ls
if [ -x /usr/bin/dircolors ]; then
    test -r ~/.dircolors && eval "$(dircolors -b ~/.dircolors)" || eval "$(dircolors -b)"
    alias ls='ls --color=auto'
    alias grep='grep --color=auto'
    alias fgrep='fgrep --color=auto'
    alias egrep='egrep --color=auto'
fi

# Useful ls aliases
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'

# Alert alias for long running commands
alias alert='notify-send --urgency=low -i "$([ $? = 0 ] && echo terminal || echo error)" "$(history|tail -n1|sed -e '\''s/^\s*[0-9]\+\s*//;s/[;&|]\s*alert$//'\'')"'

# Git shortcuts
alias g='git'
alias ga='git add .'
alias gc='git commit'

# Tmux shortcuts
alias tk='tmux kill-server'
alias tks='tmux kill-session -t'
alias tn='tmux new'
alias tns='tmux new-session -s'

# Utility aliases
alias ep='echo "$PATH" | tr ":" "\n"'

# kubectl
command -v kubectl >/dev/null 2>&1 && alias k='kubectl'
