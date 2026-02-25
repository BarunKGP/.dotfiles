# Shell environment: common variables and OS detection

# Detect operating system
_detect_os() {
    if [ -f /proc/version ] && grep -qi microsoft /proc/version 2>/dev/null; then
        echo "wsl2"
    elif [ "$(uname -s)" = "Darwin" ]; then
        echo "macos"
    else
        echo "linux"
    fi
}

export OS_TYPE="${OS_TYPE:-$(_detect_os)}"
unset -f _detect_os

# XDG Base Directory Specification
export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"

# Go paths
export GOROOT="/usr/local/go"
export GOPATH="$HOME/go"
