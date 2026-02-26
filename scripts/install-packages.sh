#!/usr/bin/env bash
# scripts/install-packages.sh - Configuration-driven package installation
# Parses packages.conf and installs packages using the detected package manager
# Can be sourced by install.sh or run standalone

set -u

# Determine if this script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    # Running standalone - initialize defaults
    DOTFILES="${DOTFILES:-.}"
    MINIMAL="${MINIMAL:-false}"
    SKIP_OPTIONAL="${SKIP_OPTIONAL:-false}"
fi

# Package manager detection
detect_package_manager() {
    if command -v brew >/dev/null 2>&1; then
        echo "brew"
    elif command -v apt-get >/dev/null 2>&1; then
        echo "apt"
    elif command -v apk >/dev/null 2>&1; then
        echo "apk"
    elif command -v dnf >/dev/null 2>&1; then
        echo "dnf"
    else
        echo "unknown"
    fi
}

# Check if package is already installed
is_installed() {
    command -v "$1" >/dev/null 2>&1
}

# Get the package name for the current package manager
get_package_name() {
    local canonical="$1"
    local tokens=("${@:2}")
    local pm="$PACKAGE_MANAGER"

    # Look for pm=name override in tokens
    for token in "${tokens[@]}"; do
        if [[ "$token" == "$pm="* ]]; then
            echo "${token#${pm}=}"
            return
        fi
    done

    # Fall back to canonical name
    echo "$canonical"
}

# Install with appropriate privilege escalation
install_package() {
    local pkg_name="$1"
    local canonical="$2"


    case "$PACKAGE_MANAGER" in
        brew)
            brew install "$pkg_name" || return 1
            ;;
        apt)
            if ! command -v sudo >/dev/null 2>&1 && [ "$(id -u)" -ne 0 ]; then
                echo "⚠️  [warn] Neither sudo nor root available; skipping apt install of $canonical"
                return 1
            fi
            if command -v sudo >/dev/null 2>&1; then
                sudo apt-get install -qq -y "$pkg_name" || return 1
            else
                apt-get install -qq -y "$pkg_name" || return 1
            fi
            ;;


        apk)
            if ! command -v sudo >/dev/null 2>&1 && [ "$(id -u)" -ne 0 ]; then
                echo "⚠️  [warn] Neither sudo nor root available; skipping apk install of $canonical"
                return 1
            fi
            if command -v sudo >/dev/null 2>&1; then
                sudo apk add --no-cache "$pkg_name" || return 1
            else
                apk add --no-cache "$pkg_name" || return 1
            fi
            ;;
        dnf)
            if ! command -v sudo >/dev/null 2>&1 && [ "$(id -u)" -ne 0 ]; then
                echo "⚠️  [warn] Neither sudo nor root available; skipping dnf install of $canonical"
                return 1
            fi
            if command -v sudo >/dev/null 2>&1; then
                sudo dnf install -y "$pkg_name" || return 1
            else
                dnf install -y "$pkg_name" || return 1
            fi
            ;;
        unknown)
            echo "⚠️  [warn] Unknown package manager for $canonical"
            return 1
            ;;
    esac
}

# Custom installer handler
run_custom_install() {
    local canonical="$1"
    local script="$DOTFILES/scripts/packages/$canonical.sh"

    if [ ! -f "$script" ]; then
        echo "⚠️  [warn] Custom install script not found: $script"
        return 1
    fi

    if [ ! -x "$script" ]; then
        chmod +x "$script"
    fi

    # Source custom script and call install function if available
    (
        export DOTFILES PACKAGE_MANAGER OS_TYPE
        . "$script"
        if declare -f install_"${canonical//-/_}" >/dev/null; then
            install_"${canonical//-/_}"
        fi
    )
}

# Main package installation function
install_packages() {
    local conf_file="$DOTFILES/packages.conf"

    if [ ! -f "$conf_file" ]; then
        echo "⚠️  [warn] packages.conf not found at $conf_file"
        return
    fi

    PACKAGE_MANAGER=$(detect_package_manager)
    echo "📦 Detected package manager: $PACKAGE_MANAGER"

    # Detect OS type (for custom installers)
    detect_os() {
        if [ -f /proc/version ] && grep -qi microsoft /proc/version 2>/dev/null; then
            echo "wsl2"
        elif [ "$(uname -s)" = "Darwin" ]; then
            echo "macos"
        else
            echo "linux"
        fi
    }
    OS_TYPE=$(detect_os)

    # Run apt-get update once if using apt
    if [ "$PACKAGE_MANAGER" = "apt" ]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo apt-get update -qq || true
        elif [ "$(id -u)" -eq 0 ]; then
            apt-get update -qq || true
        fi
    fi

    local line_num=0
    while IFS= read -r line; do
        ((++line_num))

        # Skip comments and blank lines
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "$line" ]] && continue
        [[ "$line" =~ ^[[:space:]]*$ ]] && continue

        # Parse line: canonical-name group [tokens...]
        read -r canonical group tokens <<< "$line"

        # Validate line format
        if [ -z "$canonical" ] || [ -z "$group" ]; then
            continue
        fi

        # Check if package should be skipped based on group and flags
        case "$group" in
            editor|util)
                # editor and util groups always installed
                ;;
            shell|terminal)
                if [ "$MINIMAL" = true ]; then
                    echo "⊘  [skip] $canonical (group: $group, skipped with --minimal)"
                    continue
                fi
                ;;
            optional)
                if [ "$SKIP_OPTIONAL" = true ]; then
                    echo "⊘  [skip] $canonical (group: $group, skipped with --skip-optional)"
                    continue
                fi
                # Prompt user for optional packages (only if stdin is a TTY)
                if [ -t 0 ]; then
                    read -p "Install $canonical? [y/N] " -n 1 -r response
                    echo
                    if [[ ! $response =~ ^[Yy]$ ]]; then
                        echo "⊘  [skip] $canonical (declined by user)"
                        continue
                    fi
                else
                    # Non-interactive environment: skip optional packages
                    echo "⊘  [skip] $canonical (non-interactive mode; use --skip-optional or pass --help for options)"
                    continue
                fi
                ;;
            *)
                echo "⚠️  [warn] Unknown group '$group' for package $canonical at line $line_num"
                continue
                ;;
        esac

        # Check if already installed
        if is_installed "$canonical"; then
            echo "⊘  [skip] $canonical already installed"
            continue
        fi

        # Handle custom installs
        if [[ "$tokens" =~ install=custom ]]; then
            echo "📥 Installing $canonical (custom)..."
            if run_custom_install "$canonical"; then
                echo "✓  [done] $canonical"
            else
                echo "✗  [fail] $canonical (custom install failed, continuing...)"
            fi
            continue
        fi

        # Standard package install
        echo "📥 Installing $canonical..."
        local pkg_name
        pkg_name=$(get_package_name "$canonical" $tokens)

        if install_package "$pkg_name" "$canonical"; then
            echo "✓  [done] $canonical"
        else
            echo "✗  [fail] $canonical (install failed, continuing...)"
        fi
    done < "$conf_file"

    echo "✓  Package installation complete"
}

# Only run main function if script is executed directly (not sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_packages
fi
