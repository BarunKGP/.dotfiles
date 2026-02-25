# Dotfile management and organization

This repo contains configurations for bash, zsh, tmux, neovim, git, and espanso, designed to be portable across WSL2, macOS, and Linux environments.

## Quick Start

Clone this repo on a new machine:
```sh
git clone https://github.com/BarunKGP/dotfiles.git ~/.dotfiles
cd ~/.dotfiles
./install.sh
```

The `install.sh` script will:
- Detect your OS (WSL2, macOS, or Linux)
- Check for [GNU stow](https://www.gnu.org/software/stow/) (required)
- Symlink all packages to your home directory
- Create a machine-local config file for environment-specific settings

After installation, edit `~/.config/shell/local.sh` to add machine-specific configuration (API keys, Windows paths for WSL2, etc.).

## Architecture

### Shell Configuration (Bash & Zsh)

The shell configuration is modularized and portable across all supported platforms:

- **`bash/.bashrc`** and **`zsh/.zshrc`**: Thin loaders that source modules
- **`shell/common/`**: Cross-platform modules (aliases, functions, path, tools)
- **`shell/bash/`** and **`shell/zsh/`**: Shell-specific options and completions
- **`shell/common/os/`**: OS-specific configuration (wsl2.sh, macos.sh, linux.sh)
- **`~/.config/shell/local.sh`**: Machine-local overrides (not tracked in git)

**Key Features:**
- Automatic OS detection (WSL2, macOS, Linux)
- Shared aliases and functions between bash and zsh
- Environment-specific settings loaded from `~/.config/shell/local.sh`
- Support for multiple dev environments (devcontainers, remote servers, local machines)

### Other Configurations

This repo also contains configurations for:

- **tmux**: Terminal multiplexer with plugins and theming
- **neovim**: Modern text editor with LSP support
- **git**: Global git configuration and aliases
- **espanso**: Text expansion snippets (optional)

## Setup

### Standard Installation

```sh
./install.sh
```

For minimal setup (useful in devcontainers or minimal environments):
```sh
./install.sh --minimal
```

### Manual Installation

If you prefer to stow packages individually:

```sh
stow bash zsh nvim tmux git -d /path/to/.dotfiles -t $HOME
```

### Machine-Local Configuration

After installation, edit `~/.config/shell/local.sh`:

- **API keys**: `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`
- **WSL2-specific**: `WINDOWS_USER`, `ESPANSO_PATH`, `OBSIDIAN_VAULT`, `NVIM_PATH`, `GECKODRIVER_PATH`
- **Custom aliases/functions**: Add any machine-specific shell customizations

### Shell-Specific Setup

#### Bash

```bash
source ~/.bashrc
```

#### Zsh

```zsh
source ~/.zshrc
```

### tmux

**tmux** is a lightweight, extensible terminal multiplexer that easily allows you to switch between diferent programs in one terminal.
It also has a rich plugin ecosystem and can be configured to auto-save and restore running sessions across reboots.
These features make **tmux** an important tool in any terminal-based developer workflow.

Install **tmux** according to the [instructions for your system](https://github.com/tmux/tmux/wiki/Installing) on your system.

#### Neovim

Neovim is configured via kickstart.nvim (Lua-based). The config is automatically stowed to `~/.config/nvim`.

To install Neovim:
- **macOS**: `brew install neovim` or use `setup_nvim.sh` for the appimage
- **Linux/WSL2**: Use `setup_nvim.sh` to download the portable appimage, or install via your package manager
- **devcontainers**: `apt install neovim` or use the appimage

Plugin management uses `lazy.nvim`. After starting Neovim, plugins are auto-installed.

Optional: Run `scripts/setup_nvim.sh` to download the `nvim.appimage` (portable version).

#### Git

Global git configuration is stowed to `~/.gitconfig`. Set your user identity once:

```sh
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

#### Espanso (optional)

espanso is a text expander for frequently used snippets. The tracked files only contain templates; actual sensitive data is gitignored.

**Setup:**
1. Install [espanso](https://espanso.org/docs/get-started/) for your system
2. Copy espanso config files to the espanso config directory:
   ```sh
   cp espanso/config/default.yml ~/.config/espanso/  # or Windows AppData equivalent
   ```
3. Set up match files from templates:
   ```sh
   cp espanso/match/base.example.yml ~/.config/espanso/match/base.yml
   cp espanso/match/emory.example.yml ~/.config/espanso/match/emory.yml  # if needed
   ```
4. Edit the match files to add your actual values (API keys, email addresses, paths)
5. Restart espanso: `espanso restart`

> [!WARNING]
>
> The `espanso/match/base.yml` and `espanso/match/emory.yml` files are **gitignored** for security.
> Never commit passwords, credentials, or personal information. The `.example.yml` files serve as templates.

---

## OS-Specific Notes

### WSL2

- SSH agent is automatically configured to bridge Windows ↔ WSL credentials
- Configure Windows paths in `~/.config/shell/local.sh`:
  ```bash
  export WINDOWS_USER="/mnt/c/Users/YourUsername"
  export ESPANSO_PATH="$WINDOWS_USER/AppData/Roaming/espanso/dawstored"
  export OBSIDIAN_VAULT="$WINDOWS_USER/Documents/Obsidian/obsidian-vaults/VaultName"
  ```

### macOS

- Homebrew is auto-detected and configured
- zsh is the default shell; both bash and zsh are supported

### Linux / Devcontainers

- Minimal setup with `./install.sh --minimal` installs only shell and git configs
- Full setup can be run in devcontainers but may not need tmux/nvim

---

## Troubleshooting

**Shell not loading modules:**
- Check that `$DOTFILES` is correctly set: `echo $DOTFILES`
- Verify `~/.config/shell/local.sh` exists and is readable
- Test sourcing manually: `. ~/.bashrc`

**Stow conflicts:**
- Check for existing symlinks: `ls -la ~/ | grep "\->"`
- Remove old symlinks before stowing: `rm ~/.bashrc ~/.zshrc ~/.tmux.conf`
- Use `stow --restow` to replace existing links

**macOS zsh issues:**
- macOS Catalina+ defaults to zsh; ensure `~/.zshrc` is sourced
- Add `~/.zshrc` to zsh startup if needed: `echo 'source ~/.zshrc' >> ~/.zprofile`

---

## GNU Stow

This repo uses [GNU stow](https://www.gnu.org/software/stow/) for managing dotfiles. Each package
(bash, zsh, nvim, tmux, git) is structured as a stow directory and gets symlinked to your home
directory on installation.
