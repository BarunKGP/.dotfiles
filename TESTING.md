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

### 3. Clone and Install Dotfiles

Inside the container:

```bash
# Clone the feat/modularize branch (or main once merged)
git clone --branch feat/modularize https://github.com/BarunKGP/.dotfiles.git ~/.dotfiles

# Navigate to dotfiles
cd ~/.dotfiles

# Run the installer
./install.sh
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

## Full Testing Checklist

- [ ] Container starts and SSH is accessible
- [ ] Can clone dotfiles from feat/modularize branch
- [ ] `./install.sh` runs without errors
- [ ] Both bash and zsh detect as available
- [ ] Stow creates symlinks for bash/, zsh/, git/
- [ ] `~/.bashrc` and `~/.zshrc` source correctly
- [ ] `XDG_CONFIG_HOME` is set to `~/.config`
- [ ] Machine-local config file created at `$XDG_CONFIG_HOME/shell/local.sh`
- [ ] Aliases work in bash (e.g., `ll` for `ls -alF`)
- [ ] Aliases work in zsh
- [ ] Can switch between bash and zsh
- [ ] Custom aliases in local.sh load correctly
- [ ] Git config is symlinked
