# Testing Dotfiles with Docker (Alpine Linux)

This guide explains how to test the dotfiles installation on Alpine Linux using Docker.

## Prerequisites

- Docker installed
- Docker Compose installed
- SSH client

## Quick Start

### 1. Build and Start the Container

```bash
docker-compose up -d
```

This will:
- Build the Alpine Linux image with all prerequisites (git, stow, bash, zsh, openssh)
- Start the SSH server on `localhost:2222`
- Create a `testuser` account for testing

### 2. SSH into the Container

```bash
ssh -p 2222 testuser@localhost
```

When prompted for password, use any password (the container allows password login for testing).

### 3. Dotfiles Already Installed

Dotfiles are automatically installed during container build using `dotctl install --skip-optional`.

To verify the installation is complete:

```bash
# Check that dotfiles are installed
ls -la ~/.bashrc ~/.zshrc ~/.gitconfig  # Should be symlinks

# Verify dotctl is in PATH
which dotctl

# Check local config was created
ls -la $XDG_CONFIG_HOME/shell/local.sh
```

To run dotctl again (e.g., with different flags):

```bash
cd ~/.dotfiles
dotctl install --help  # Show all options
```

### 4. Verify Installation

After installation, you can verify:

```bash
# Check bash configuration
bash -i -c 'echo $OS_TYPE'  # Should output: linux
bash -i -c 'echo $XDG_CONFIG_HOME'  # Should output: /home/testuser/.config

# Check zsh configuration
zsh -i -c 'echo $OS_TYPE'  # Should output: linux
zsh -i -c 'echo $XDG_CONFIG_HOME'  # Should output: /home/testuser/.config

# Check stowed files
ls -la ~/.bashrc ~/.zshrc ~/.gitconfig  # Should be symlinks

# Test shell aliases
source ~/.bashrc
alias | grep "ll="  # Should show: alias ll='ls -alF'

# Verify tmux config
ls -la ~/.tmux.conf  # Should be a symlink
```

### 5. Test Shell Switching

Switch between shells to ensure both work:

```bash
# Current shell (should be bash by default)
echo $SHELL

# Switch to zsh
zsh
echo $OS_TYPE  # Should output: linux

# Back to bash
bash
```

## Testing Specific Features

### Test XDG_CONFIG_HOME

```bash
echo $XDG_CONFIG_HOME  # Should be /home/testuser/.config
ls $XDG_CONFIG_HOME/shell/local.sh  # Should exist and be readable
```

### Test Machine-Local Config

Edit the local config file:

```bash
nano $XDG_CONFIG_HOME/shell/local.sh
```

Add a custom alias:

```bash
alias mytest='echo "Custom alias works!"'
```

Reload shell and test:

```bash
exec $SHELL
mytest  # Should output: Custom alias works!
```

### Test OS Detection

Verify OS detection is working:

```bash
bash -i -c 'echo "OS: $OS_TYPE"'  # Should output: OS: linux
bash -i -c 'echo "GOPATH: $GOPATH"'  # Should output: GO path
```

### Test Available Tools

```bash
# Check if NVM is available
nvm --version  # May not be installed, that's ok

# Check if pyenv is available
pyenv --version  # May not be installed, that's ok

# Check zoxide
zoxide --version
```

## Container Management

### View Logs

```bash
docker-compose logs -f
```

### Stop the Container

```bash
docker-compose down
```

### Remove Everything

```bash
docker-compose down -v  # Removes volumes too
```

### Rebuild Image

If you make changes to the Dockerfile:

```bash
docker-compose up -d --build
```

### Access Container Shell Directly

Without SSH:

```bash
docker-compose exec -u testuser alpine-dotfiles-test /bin/bash
```

## Troubleshooting

### SSH Connection Refused

The container may need a moment to start SSH. Try:

```bash
docker-compose restart
sleep 2
ssh -p 2222 testuser@localhost
```

### Permission Denied

Make sure you're using the correct port:

```bash
ssh -p 2222 testuser@localhost  # Correct
ssh testuser@localhost  # Wrong - uses port 22
```

### Stow Installation Fails

Verify stow is available in the container:

```bash
docker-compose exec -u testuser alpine-dotfiles-test stow --version
```

### Shell Not Loading Config

Check the file locations:

```bash
ls -la ~/.bashrc ~/.zshrc
cat ~/.bashrc | head -5
```

## Notes

- **Alpine Linux** is a minimal Linux distribution (~5MB) - good for testing
- **Package manager**: Alpine uses `apk` (not apt or yum)
- **Init system**: Alpine doesn't use systemd by default
- **SSH**: Configured for password authentication for easy testing
- **Home directory**: `/home/testuser` - shared via Docker volume for persistence

## Full Testing Checklist (e2e Verification)

After `ssh -p 2222 testuser@localhost`:

**Environment Setup**
- [ ] `$OS_TYPE` = linux (test: `bash -i -c 'echo $OS_TYPE'`)
- [ ] `$XDG_CONFIG_HOME` = /home/testuser/.config (test: `echo $XDG_CONFIG_HOME`)
- [ ] `$HOME` = /home/testuser (test: `echo $HOME`)

**Symlinks (Stow)**
- [ ] `~/.bashrc` is a symlink to `~/.dotfiles/bash/.bashrc`
- [ ] `~/.zshrc` is a symlink to `~/.dotfiles/zsh/.zshrc`
- [ ] `~/.gitconfig` is a symlink to `~/.dotfiles/git/.gitconfig`
- [ ] `~/.tmux.conf` is a symlink to `~/.dotfiles/tmux/.tmux.conf`

**Shell Config & Aliases**
- [ ] Aliases work in bash: `ll` expands to `ls -alF` (test: `bash -i -c 'alias | grep ll'`)
- [ ] Aliases work in zsh: `ll` expands to `ls -alF` (test: `zsh -i -c 'alias | grep ll'`)

**Local Configuration**
- [ ] `$XDG_CONFIG_HOME/shell/local.sh` exists and is readable
- [ ] Custom aliases in `local.sh` load after reload (add `alias mytest='echo OK'`, then `exec $SHELL && mytest`)

**Idempotency**
- [ ] Second `dotctl install` run succeeds without errors
- [ ] No conflicts or symlink errors on re-run

**Binary & Tools**
- [ ] `dotctl --version` or `which dotctl` works
- [ ] `stow --version` works
- [ ] `zsh --version` works
