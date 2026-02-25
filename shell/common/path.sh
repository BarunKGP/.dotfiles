# PATH management - add standard tool directories

# Go binaries
export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"

# Local user bin
[ -d "$HOME/.local/bin" ] && export PATH="$HOME/.local/bin:$PATH"

# Remove duplicate entries from PATH
# This must be run last after all PATH additions
export PATH="$(echo "$PATH" | awk -v RS=: '!a[$1]++' | paste -sd: -)"
