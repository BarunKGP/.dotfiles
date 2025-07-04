:'
This script enables sharing SSH credentials between Windwows and WSL
when used on the same machine, standardizing things such as Git access
across platforms.

NOTE: Do not use this for sharing SSH keys unless they belong to the same
user on the same machine.
'

#!/bin/bash

# Location to store agent environment
AGENT_ENV="$HOME/.ssh/agent.env"

# Function to start ssh-agent
start_agent() {
    echo "Starting a new ssh-agent..."
    eval "$(ssh-agent -s)" >/dev/null
    ssh-add ~/.ssh/git_ed25519
    # Save environment variables
    echo "export SSH_AUTH_SOCK=$SSH_AUTH_SOCK" >"$AGENT_ENV"
    echo "export SSH_AGENT_PID=$SSH_AGENT_PID" >>"$AGENT_ENV"
    chmod 600 "$AGENT_ENV"
}

# Check if environment file exists
if [ -f "$AGENT_ENV" ]; then
    source "$AGENT_ENV" >/dev/null
    # Check if agent is actually running
    if ! kill -0 "$SSH_AGENT_PID" 2>/dev/null; then
        start_agent
    fi
else
    start_agent
fi
