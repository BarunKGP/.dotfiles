# macOS-specific shell configuration

# Homebrew initialization
if [ -x /opt/homebrew/bin/brew ]; then
    # Apple Silicon Macs
    eval "$(/opt/homebrew/bin/brew shellenv)"
elif [ -x /usr/local/bin/brew ]; then
    # Intel Macs
    eval "$(/usr/local/bin/brew shellenv)"
fi
