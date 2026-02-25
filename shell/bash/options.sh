# Bash-specific shell options and history configuration

# Don't put duplicate lines or lines starting with space in history
HISTCONTROL=ignoreboth

# Append to history file instead of overwriting
shopt -s histappend

# Set history length
HISTSIZE=1000
HISTFILESIZE=2000

# Update window size after each command
shopt -s checkwinsize
