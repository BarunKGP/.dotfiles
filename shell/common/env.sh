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

# Source generated language environment (created by dotctl env render)
if [ -f "$XDG_CONFIG_HOME/shell/generated/languages.sh" ]; then
	. "$XDG_CONFIG_HOME/shell/generated/languages.sh"
fi
