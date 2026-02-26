FROM alpine:latest

# Install required packages for bootstrapping dotfiles
# zsh, neovim, tmux are now installed by dotfiles via packages.conf
RUN apk add --no-cache \
    openssh \
    openssh-client \
    git \
    stow \
    bash \
    curl \
    vim \
    util-linux \
    sudo

# Create a non-root user for testing with password "test" and bash as default shell
RUN addgroup -g 1000 testuser && \
    adduser -D -u 1000 -G testuser -s /bin/bash testuser && \
    echo "testuser:test" | chpasswd

# Grant testuser passwordless sudo for apk (minimal privilege for package installation)
RUN mkdir -p /etc/sudoers.d && \
    echo "testuser ALL=(root) NOPASSWD: /sbin/apk" >> /etc/sudoers.d/testuser && \
    chmod 0440 /etc/sudoers.d/testuser

# Set up SSH
RUN ssh-keygen -A && \
    mkdir -p /home/testuser/.ssh && \
    chmod 700 /home/testuser/.ssh

# Allow SSH login with password (for easy testing)
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config && \
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/^PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config

# Create entrypoint script
RUN mkdir -p /home/testuser && chown -R testuser:testuser /home/testuser

# Display helpful info
RUN cat > /etc/motd << 'EOF'
Welcome to Alpine Linux Dotfiles Test Container!

To test the dotfiles installation:
1. Clone the repository: git clone https://github.com/BarunKGP/.dotfiles.git ~/.dotfiles
2. Or: git clone --branch feat/modularize https://github.com/BarunKGP/.dotfiles.git ~/.dotfiles
3. Navigate: cd ~/.dotfiles
4. Run installer: ./install.sh
5. Reload shell: exec $SHELL

Installation flags:
  ./install.sh --minimal      # Skip nvim/tmux (fast for devcontainers)
  ./install.sh --no-packages  # Manual package management, only stow dotfiles
  ./install.sh --skip-optional # Skip optional packages like espanso

Note: zsh, neovim, and tmux are now installed by dotfiles via packages.conf

SSH Credentials: testuser / test
EOF

# Start SSH server
CMD ["/usr/sbin/sshd", "-D", "-e"]

EXPOSE 22
