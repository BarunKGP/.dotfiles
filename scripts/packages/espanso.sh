#!/usr/bin/env bash
# scripts/packages/espanso.sh - Custom platform-aware espanso installer
# Called by scripts/install-packages.sh when install=custom is detected
# Handles OS-specific installation methods per https://espanso.org/docs/install/linux/

# Install espanso based on OS and package manager
install_espanso() {
    local os="${OS_TYPE:-linux}"
    local pm="${PACKAGE_MANAGER:-unknown}"

    echo "Installing espanso on $os with $pm..."

    case "$os" in
        wsl2)
            echo ""
            echo "╔════════════════════════════════════════════════════════════╗"
            echo "║  Espanso on WSL2 - Windows Installation Required          ║"
            echo "╚════════════════════════════════════════════════════════════╝"
            echo ""
            echo "Espanso is a system-level text expander that must run on Windows"
            echo "to intercept keystrokes. It cannot run inside WSL2."
            echo ""
            echo "To install espanso on your Windows host:"
            echo "  1. Open PowerShell as Administrator"
            echo "  2. Run: winget install Espanso.Espanso"
            echo ""
            echo "For more details: https://espanso.org/docs/install/windows/"
            echo ""
            return 0
            ;;

        macos)
            echo "Installing espanso via Homebrew..."
            if command -v brew >/dev/null 2>&1; then
                brew install espanso
                echo "✓ Espanso installed successfully"
                return 0
            else
                echo "✗ Homebrew not found; cannot install espanso"
                return 1
            fi
            ;;

        linux)
            case "$pm" in
                apt)
                    _install_espanso_deb
                    return $?
                    ;;
                apk)
                    echo "⚠️  Alpine Linux: Espanso is not available in Alpine repositories"
                    echo "  Consider using a different distribution or building from source"
                    return 1
                    ;;
                dnf)
                    _install_espanso_fedora
                    return $?
                    ;;
                *)
                    echo "⚠️  Unsupported package manager ($pm) for Linux espanso install"
                    return 1
                    ;;
            esac
            ;;

        *)
            echo "⚠️  Unknown OS type: $os"
            return 1
            ;;
    esac
}

# Download and install espanso .deb for Ubuntu/Debian
_install_espanso_deb() {
    local tmpdir
    tmpdir=$(mktemp -d)
    trap "rm -rf $tmpdir" EXIT

    echo "Downloading espanso .deb package..."

    # Get latest release from GitHub API
    local latest_release
    latest_release=$(curl -s "https://api.github.com/repos/espanso/espanso/releases/latest" | grep '"tag_name"' | head -1 | sed -E 's/.*"v([^"]+)".*/\1/')

    if [ -z "$latest_release" ]; then
        echo "✗ Could not determine latest espanso version"
        return 1
    fi

    local deb_url="https://github.com/espanso/espanso/releases/download/v${latest_release}/espanso-debian-x11-${latest_release}-amd64.deb"

    if ! curl -fsSL -o "$tmpdir/espanso.deb" "$deb_url"; then
        echo "✗ Failed to download espanso from $deb_url"
        return 1
    fi

    echo "Installing espanso .deb..."
    if command -v sudo >/dev/null 2>&1; then
        sudo apt install -y "$tmpdir/espanso.deb"
    else
        apt install -y "$tmpdir/espanso.deb"
    fi

    if command -v espanso >/dev/null 2>&1; then
        echo "✓ Espanso installed successfully"
        return 0
    else
        echo "✗ Espanso installation verification failed"
        return 1
    fi
}

# Install espanso on Fedora via dnf
_install_espanso_fedora() {
    echo "Installing espanso on Fedora..."
    echo "Note: Requires the Terra repo to be configured"
    echo ""

    # Check if Terra repo is configured
    if ! dnf repolist | grep -q terra; then
        echo "⚠️  Terra repository not found. Setting up..."
        if command -v sudo >/dev/null 2>&1; then
            sudo dnf copr enable espanso/espanso
        else
            dnf copr enable espanso/espanso
        fi
    fi

    if command -v sudo >/dev/null 2>&1; then
        sudo dnf install -y espanso-x11
    else
        dnf install -y espanso-x11
    fi

    if command -v espanso >/dev/null 2>&1; then
        echo "✓ Espanso installed successfully"
        return 0
    else
        echo "✗ Espanso installation verification failed"
        return 1
    fi
}

# Only run if script is executed directly or sourced and called
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_espanso
fi
