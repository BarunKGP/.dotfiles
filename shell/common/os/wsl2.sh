# WSL2-specific shell configuration

# SSH agent for Windows ↔ WSL credential sharing
. "$DOTFILES/scripts/startup/ssh_agent_setup.sh"

# Windows integration configuration
# The following environment variables should be set in ~/.config/shell/local.sh:
#
#   export WINDOWS_USER="/mnt/c/Users/YourUsername"
#   export ESPANSO_PATH="$WINDOWS_USER/AppData/Roaming/espanso/dawstored"
#   export OBSIDIAN_VAULT="$WINDOWS_USER/Documents/Obsidian/obsidian-vaults/VaultName"
#   export NVIM_PATH="/opt/nvim-linux-x86_64/bin/"
#   export GECKODRIVER_PATH="$HOME/.config/geckodriver"
#
# Once set, they will be added to PATH automatically below:

[ -n "$NVIM_PATH" ] && export PATH="$PATH:$NVIM_PATH"
[ -n "$GECKODRIVER_PATH" ] && export PATH="$PATH:$GECKODRIVER_PATH"
