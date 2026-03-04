#!/usr/bin/env bash
# scripts/packages/pyenv.sh - Custom pyenv (Python Version Manager) installer
# Called by scripts/install-packages.sh when install=custom is detected
# Idempotent installer with dry-run support

# Install pyenv for multi-version Python management
install_pyenv() {
    # Check for dry-run mode
    if [ "${DRY_RUN:-0}" = "1" ]; then
        echo "Would install pyenv"
        return 0
    fi

    # Check if pyenv is already installed (idempotency)
    if [ -f "$HOME/.pyenv/bin/pyenv" ]; then
        echo "✓ pyenv already installed at $HOME/.pyenv"
        return 0
    fi

    echo "Installing pyenv..."

    # Download and run the pyenv installer
    if ! curl https://pyenv.run | bash; then
        echo "✗ Failed to install pyenv"
        return 1
    fi

    echo "✓ pyenv installed successfully"
    echo "Note: pyenv requires shell configuration. Run: dotctl env render"
    return 0
}

# Only run if script is executed directly or sourced and called
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_pyenv
fi
