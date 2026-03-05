# Stage 1: Build dotctl binary from standalone dotctl repository
FROM golang:1.22-alpine AS builder

ARG DOTCTL_REPO=https://github.com/BarunKGP/dotctl.git
ARG DOTCTL_BRANCH=main

RUN apk add --no-cache git

RUN git clone --branch ${DOTCTL_BRANCH} ${DOTCTL_REPO} /build
WORKDIR /build
RUN go build -o bin/dotctl ./cmd/dotctl

# Stage 2: Runtime (Alpine, no Go)
FROM alpine:latest

ARG DOTFILES_REPO=https://github.com/BarunKGP/.dotfiles.git
ARG DOTFILES_BRANCH=feat/modularize
ARG INSTALL_FLAGS=--skip-optional

# Bootstrap: only what's needed before dotctl install runs
RUN apk add --no-cache \
    openssh \
    openssh-client \
    git \
    stow \
    bash \
    curl \
    vim \
    util-linux \
    sudo \
    shadow

# Copy dotctl binary from builder stage
COPY --from=builder /build/bin/dotctl /usr/local/bin/dotctl

# Create testuser (bash as initial shell — zsh installed by dotfiles)
RUN addgroup -g 1000 testuser && \
    adduser -D -u 1000 -G testuser -s /bin/bash testuser && \
    echo "testuser:test" | chpasswd

# Passwordless sudo for apk (package install) and usermod (shell change)
RUN mkdir -p /etc/sudoers.d && \
    echo "testuser ALL=(root) NOPASSWD: /sbin/apk, /usr/sbin/usermod" \
        >> /etc/sudoers.d/testuser && \
    chmod 0440 /etc/sudoers.d/testuser

# Clone dotfiles from DOTFILES_BRANCH (default feat/modularize with shell config)
RUN mkdir -p /home/testuser && chown -R testuser:testuser /home/testuser
USER testuser
ENV HOME=/home/testuser
RUN git clone --branch ${DOTFILES_BRANCH} ${DOTFILES_REPO} ${HOME}/.dotfiles

# Run dotctl install (as root so it can access apk in /sbin, but install for testuser)
USER root
RUN cd /home/testuser/.dotfiles && HOME=/home/testuser dotctl install ${INSTALL_FLAGS} && \
    chown -R testuser:testuser /home/testuser
USER testuser

# Set zsh as default shell now that it's installed (bypasses PAM, no chsh needed)
USER root
RUN if command -v zsh >/dev/null 2>&1; then \
        usermod -s "$(command -v zsh)" testuser; \
    fi

# SSH setup
RUN ssh-keygen -A && \
    mkdir -p /home/testuser/.ssh && \
    chmod 700 /home/testuser/.ssh && \
    chown -R testuser:testuser /home/testuser/.ssh

# Enable password auth (Alpine sshd_config defaults allow it, but be explicit)
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config && \
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
    echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config

RUN cat > /etc/motd << 'EOF'
Alpine Linux Dotfiles Container — ready to code.

SSH credentials: testuser / test
Default shell:   zsh (with dotfiles pre-installed)

To rebuild fresh:
  docker-compose down && docker-compose up -d --build
EOF

CMD ["/usr/sbin/sshd", "-D", "-e"]
EXPOSE 22
