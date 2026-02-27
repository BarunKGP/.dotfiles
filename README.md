# Dotfiles — Portable Configuration for Bash, Zsh, Tmux, Neovim & Git

This repo contains portable shell, editor, and terminal configurations for bash, zsh, tmux, neovim, git, and espanso—designed to work seamlessly across **WSL2, macOS, and Linux**.

## Quick Start

Clone this repo and run the installer:

```sh
git clone https://github.com/BarunKGP/dotfiles.git ~/.dotfiles
cd ~/.dotfiles
dotctl install
```

The installer will:
- ✅ Detect your OS (WSL2, macOS, or Linux)
- ✅ Detect your package manager (apt, apk, brew, dnf)
- ✅ Install packages (zsh, neovim, tmux, fzf, ripgrep, espanso, etc.)
- ✅ Stow dotfiles to symlink configs to your home directory
- ✅ Set zsh as default shell (if available)
- ✅ Create `~/.config/shell/local.sh` for machine-specific settings

After installation, customize `~/.config/shell/local.sh` with your personal configuration (API keys, Windows paths for WSL2, etc.).

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

## Installation

### Using dotctl (Recommended)

`dotctl` is a Go-based CLI tool that provides a modern, testable installer with feature parity to the original shell scripts.

```sh
dotctl install
```

**Flags:**
- `--minimal` — Skip editor (nvim) and terminal (tmux) packages
- `--no-packages` — Skip package installation (stow configs only)
- `--skip-optional` — Skip optional packages (espanso)
- `--dry-run` — Preview what would be installed without making changes
- `--non-interactive` — Skip interactive prompts

**Examples:**
```sh
# Full install with all packages
dotctl install

# Minimal setup for devcontainers
dotctl install --minimal

# Just stow configs, no package manager
dotctl install --no-packages

# Dry-run preview
dotctl install --dry-run
```

### Alternative: Shell-Based Installation (Legacy)

For compatibility, the original shell installer is still available:

```sh
./install.sh
```

For minimal setup:
```sh
./install.sh --minimal
```

### Manual Installation

If you prefer to stow packages individually without the installer:

```sh
# Stow specific packages
stow bash zsh nvim tmux git -d ~/.dotfiles -t $HOME

# Or use dotctl's dry-run to see what would be stowed
dotctl stow --help
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

## dotctl — Modern Go-Based Installer

**dotctl** is a portable, testable CLI tool that replaces the original shell-based installer with:

- **OS Detection** — Automatically detects WSL2, macOS, or Linux
- **Package Manager Support** — Works with apt, apk, brew, dnf (auto-detected)
- **Stow Integration** — Symlinks dotfiles with conflict detection
- **Dry-run Mode** — Preview changes before applying: `dotctl install --dry-run`
- **Idempotent** — Safe to run multiple times (second run skips completed steps)
- **Tested** — Verified with unit tests and e2e Docker testing

### Using dotctl

```sh
# View available commands
dotctl --help

# Full installation
dotctl install

# View what would be installed
dotctl install --dry-run

# Health check
dotctl doctor
```

### Testing

The installation is tested using Docker with Alpine Linux. To run the e2e tests:

```sh
docker-compose build
docker-compose up -d
docker-compose exec -u testuser alpine-dotfiles-test bash -c '
  echo "Testing symlinks..."
  ls -l ~/.bashrc ~/.zshrc ~/.gitconfig ~/.tmux.conf

  echo "Testing aliases..."
  bash -i -c "alias ll"

  echo "Testing idempotency..."
  cd ~/.dotfiles && dotctl install --skip-optional
'
```

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

### Installation Issues

**dotctl install fails:**
```bash
# Check system dependencies
dotctl doctor

# Preview installation (no changes made)
dotctl install --dry-run

# Run with verbose output
dotctl install --verbose
```

**Package manager not detected:**
- Ensure your package manager (apt, brew, apk, dnf) is installed
- Use `--no-packages` to skip package installation and stow configs only

### Shell Configuration Issues

**Shell not loading modules:**
- Check that `$DOTFILES` is correctly set: `echo $DOTFILES`
- Verify `~/.config/shell/local.sh` exists and is readable
- Test sourcing manually: `. ~/.bashrc`

**Stow conflicts:**
- Check for existing symlinks: `ls -la ~/ | grep "\->"`
- Remove old symlinks before stowing: `rm ~/.bashrc ~/.zshrc ~/.tmux.conf`
- Use `stow --restow` to replace existing links (or `dotctl install` which does this automatically)

**macOS zsh issues:**
- macOS Catalina+ defaults to zsh; ensure `~/.zshrc` is sourced
- Add `~/.zshrc` to zsh startup if needed: `echo 'source ~/.zshrc' >> ~/.zprofile`

---

## GNU Stow

This repo uses [GNU stow](https://www.gnu.org/software/stow/) for managing dotfiles. Each package
(bash, zsh, nvim, tmux, git) is structured as a stow directory and gets symlinked to your home
directory on installation.
