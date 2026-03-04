#!/usr/bin/env bash
# scripts/packages/nvm.sh - Custom nvm (Node Version Manager) installer
# Called by scripts/install-packages.sh when install=custom is detected
# Idempotent installer with dry-run support

# Install nvm and a default Node.js LTS version
install_nvm() {
    # Check for dry-run mode
    if [ "${DRY_RUN:-0}" = "1" ]; then
        echo "Would install nvm and Node.js LTS"
        return 0
    fi

    # Check if nvm is already installed (idempotency)
    if [ -s "$HOME/.nvm/nvm.sh" ]; then
        echo "✓ nvm already installed at $HOME/.nvm"
        return 0
    fi

    echo "Installing nvm..."

    # Download and run the nvm install script
    local nvm_version="v0.39.7"
    if ! curl -o- "https://raw.githubusercontent.com/nvm-sh/nvm/${nvm_version}/install.sh" | bash; then
        echo "✗ Failed to install nvm"
        return 1
    fi

    # Source nvm to use it in this script
    export NVM_DIR="$HOME/.nvm"
    # shellcheck source=/dev/null
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

    # Install the LTS version of Node.js
    echo "Installing Node.js LTS..."
    if ! nvm install --lts; then
        echo "⚠️  nvm installed but Node.js LTS installation failed"
        echo "  You can install Node.js manually later with: nvm install --lts"
        return 0
    fi

    echo "✓ nvm and Node.js LTS installed successfully"
    return 0
}

# Only run if script is executed directly or sourced and called
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_nvm
fi
