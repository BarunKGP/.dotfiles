#!/bin/bash
# Thin shim for backward compatibility
# Delegates to dotctl install - the new Go-based installer
# Maintained for backward compatibility with existing documentation

set -e

# Find dotctl binary (either in PATH or in ./bin/)
if command -v dotctl >/dev/null 2>&1; then
    DOTCTL="dotctl"
elif [ -f "./bin/dotctl" ]; then
    DOTCTL="./bin/dotctl"
elif [ -f "bin/dotctl" ]; then
    DOTCTL="bin/dotctl"
else
    echo "Error: dotctl not found in PATH or ./bin/" >&2
    echo "" >&2
    echo "To build dotctl, run:" >&2
    echo "  make build" >&2
    exit 1
fi

# Pass all arguments through to dotctl install
exec "$DOTCTL" install "$@"
