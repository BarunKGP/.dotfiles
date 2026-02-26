#!/usr/bin/env bash
# Dotfiles installation script
# Sets up dotfiles on a new machine with automatic OS detection

set -euo pipefail

DOTFILES="$(cd "$(dirname "$0")" && pwd)"
MINIMAL=false
NO_PACKAGES=false
SKIP_OPTIONAL=false

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --minimal)
            MINIMAL=true
            shift
            ;;
        --no-packages)
            NO_PACKAGES=true
            shift
            ;;
        --skip-optional)
            SKIP_OPTIONAL=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "OPTIONS:"
            echo "  --minimal       Skip shell/terminal packages and plugin setup (for devcontainers)"
            echo "  --no-packages   Skip package installation; only stow dotfiles"
            echo "  --skip-optional Skip optional packages (like espanso) without prompting"
            echo "  --help          Show this help message"
            echo ""
            echo "EXAMPLES:"
            echo "  $0                    # Full install with prompts"
            echo "  $0 --minimal          # Fast setup for devcontainers (skips zsh/tmux)"
            echo "  $0 --no-packages      # Manual package mgmt; only apply dotfiles"
            echo "  $0 --skip-optional    # Install core packages, skip optional ones"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# OS detection
detect_os() {
    if [ -f /proc/version ] && grep -qi microsoft /proc/version 2>/dev/null; then
        echo "wsl2"
    elif [ "$(uname -s)" = "Darwin" ]; then
        echo "macos"
    else
        echo "linux"
    fi
}

OS=$(detect_os)

# XDG Base Directory Specification
export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"

echo "================================================"
echo "  Dotfiles Installation"
echo "================================================"
echo "Location: $DOTFILES"
echo "OS: $OS"
echo "Minimal mode: $MINIMAL"
echo ""

# Ensure all critical scripts have executable permissions
# (git sometimes doesn't preserve permissions across clones/filesystems)
ensure_executable() {
    local script="$1"
    if [ -f "$script" ] && [ ! -x "$script" ]; then
        chmod +x "$script"
    fi
}

ensure_executable "$DOTFILES/scripts/install-packages.sh"
ensure_executable "$DOTFILES/scripts/packages/espanso.sh"

# Install packages from packages.conf if enabled
if [ "$NO_PACKAGES" = false ]; then
    echo "Installing packages from packages.conf..."
    export MINIMAL DOTFILES SKIP_OPTIONAL
    . "$DOTFILES/scripts/install-packages.sh"
    install_packages
    echo ""
fi

echo "Detecting available shells..."

# Check for GNU stow
if ! command -v stow >/dev/null 2>&1; then
    echo "ERROR: GNU stow is required but not installed."
    echo ""
    echo "Install stow:"
    case "$OS" in
        wsl2|linux)
            echo "  Ubuntu/Debian: sudo apt install stow"
            echo "  Fedora: sudo dnf install stow"
            ;;
        macos)
            echo "  Homebrew: brew install stow"
            ;;
    esac
    exit 1
fi

# Stow packages
stow_package() {
    local pkg="$1"
    echo "Stowing package: $pkg"
    stow "$pkg" -d "$DOTFILES" -t "$HOME" --restow
}

# Set default shell
set_default_shell() {
    local shell="$1"
    local shell_path

    # Get the full path to the shell
    shell_path=$(command -v "$shell")

    if [ -z "$shell_path" ]; then
        echo "⚠️  $shell not found; skipping default shell change"
        return 1
    fi

    # Check if it's already the default
    if [ "$SHELL" = "$shell_path" ]; then
        echo "✓ $shell is already the default shell"
        return 0
    fi

    # Try to change the default shell
    if chsh -s "$shell_path" >/dev/null 2>&1; then
        echo "✓ Changed default shell to: $shell_path"
        return 0
    else
        echo "⚠️  Could not change default shell to $shell (requires re-login to take effect)"
        echo "   To set manually: chsh -s $shell_path"
        return 1
    fi
}

# Shell detection
detect_shells() {
    local shells=()
    if command -v bash >/dev/null 2>&1; then
        shells+=("bash")
    fi
    if command -v zsh >/dev/null 2>&1; then
        shells+=("zsh")
    fi
    echo "${shells[@]}"
}

AVAILABLE_SHELLS=$(detect_shells)

# Check if at least one shell is available
if [ -z "$AVAILABLE_SHELLS" ]; then
    echo "ERROR: Neither bash nor zsh found on this system."
    echo "At least one shell is required to use these dotfiles."
    exit 1
fi

echo "Installing packages with stow..."
echo "Available shells: $AVAILABLE_SHELLS"
for shell in $AVAILABLE_SHELLS; do
    stow_package "$shell"
done

stow_package "git"
stow_package "nvim"

if [ "$MINIMAL" = false ]; then
    stow_package "zsh"
    stow_package "tmux"
fi

# Set zsh as default shell if available
if command -v zsh >/dev/null 2>&1; then
    echo ""
    set_default_shell "zsh" || true
fi

echo ""

# Create local shell config if it doesn't exist
LOCAL_CONFIG="$XDG_CONFIG_HOME/shell/local.sh"
if [ ! -f "$LOCAL_CONFIG" ]; then
    mkdir -p "$(dirname "$LOCAL_CONFIG")"

    cat > "$LOCAL_CONFIG" <<'EOF'
# Machine-local shell configuration - NOT tracked in git
# Edit this file to add machine-specific settings

# If DOTFILES is not at the default location ~/.dotfiles, set it here:
# export DOTFILES="/path/to/your/.dotfiles"

# API keys and secrets (do NOT commit these)
# export OPENAI_API_KEY="your-key-here"
# export ANTHROPIC_API_KEY="your-key-here"

# WSL2-specific configuration (uncomment if on WSL2)
# export WINDOWS_USER="/mnt/c/Users/YourUsername"
# export ESPANSO_PATH="$WINDOWS_USER/AppData/Roaming/espanso/dawstored"
# export OBSIDIAN_VAULT="$WINDOWS_USER/Documents/Obsidian/obsidian-vaults/VaultName"
# export NVIM_PATH="/opt/nvim-linux-x86_64/bin/"
# export GECKODRIVER_PATH="$HOME/.config/geckodriver"

# macOS-specific configuration (uncomment if on macOS)
# export HOMEBREW_PREFIX="/opt/homebrew"

# Other machine-specific aliases or functions:
# alias myalias='my-command'
EOF

    echo "Created local config template: $LOCAL_CONFIG"
    echo "Edit it to add machine-specific settings (API keys, Windows paths, etc.)"
    echo ""
fi

# OS-specific guidance
case "$OS" in
    wsl2)
        echo "WSL2 Detected:"
        echo "  - Edit $LOCAL_CONFIG to configure Windows paths"
        echo "  - SSH agent will use Windows credentials via ssh_agent_setup.sh"
        ;;
    macos)
        echo "macOS Detected:"
        if ! command -v brew >/dev/null 2>&1; then
            echo "  - Homebrew not found. Install from https://brew.sh"
        fi
        ;;
esac

echo ""
echo "================================================"
echo "  Installation Complete!"
echo "================================================"
echo ""
echo "Next steps:"
echo "  1. Edit $LOCAL_CONFIG for machine-specific settings"
echo "  2. Reload your shell: exec \$SHELL"
echo ""

if [ "$MINIMAL" = false ]; then
    echo "Optional setup:"
    echo "  - Neovim plugins: In nvim, run :Lazy to manage plugins"
    echo "  - Tmux plugins: Press prefix + I to install plugins"
    echo ""
fi

# Git configuration reminder
echo "Git configuration:"
echo "  Set your git user (one-time):"
echo "    git config --global user.name 'Your Name'"
echo "    git config --global user.email 'your.email@example.com'"
echo ""
