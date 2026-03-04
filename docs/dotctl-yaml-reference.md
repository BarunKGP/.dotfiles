# dotctl.yaml Configuration Reference

## Overview

`dotctl.yaml` is the primary configuration file for dotctl. It defines:
- **Profiles**: Named groups of packages for different installation scenarios
- **Packages**: Package definitions with manager-specific installation methods
- **Container**: Docker/container image configuration (Phase 006)
- **Deploy**: Deployment configuration for remote systems (Phase 007)

## File Location

Place `dotctl.yaml` in the root of your dotfiles repository. If it's missing, dotctl falls back to `packages.conf` for backward compatibility.

## Schema

### Profiles

Profiles are named collections of package groups. Use profiles to define different installation scenarios:

```yaml
profiles:
  default:       # Standard installation (used by default)
    - util
    - shell
    - editor
    - terminal

  minimal:       # Minimal setup for devcontainers
    - shell

  workstation:   # Full developer workstation
    - util
    - shell
    - editor
    - terminal
    - languages

  devcontainer:  # Devcontainer-optimized setup
    - util
    - shell
    - editor
```

**Usage:**
```bash
# Install using default profile
dotctl install

# Install using minimal profile
dotctl install --profile minimal

# Install using workstation profile
dotctl install --profile workstation
```

### Packages

Each package defines how it should be installed across different package managers:

```yaml
packages:
  zsh:
    group: shell
    install:
      apt: zsh
      apk: zsh
      brew: zsh
      dnf: zsh
```

**Fields:**
- `group`: Package group (util, shell, editor, terminal, languages, optional, etc.)
- `install`: Installation configuration
  - `apt`: Package name for apt (Debian/Ubuntu)
  - `apk`: Package name for apk (Alpine)
  - `brew`: Package name for Homebrew (macOS)
  - `dnf`: Package name for dnf (Fedora)
  - `custom`: Path to custom installation script (e.g., `scripts/packages/nvm.sh`)

**Example with Custom Install:**
```yaml
nvm:
  group: languages
  install:
    custom: scripts/packages/nvm.sh
```

### Package Manager Name Resolution

When the same package name works across all managers, you don't need overrides:

```yaml
python3:
  group: languages
  install:
    apt: python3
    apk: python3
    brew: python3
    dnf: python3
```

When names differ (e.g., apt uses `golang-go` but others use `go`):

```yaml
go:
  group: languages
  install:
    apt: golang-go    # Special case for apt
    apk: go
    brew: go
    dnf: golang
```

### Container Configuration

Configuration for building development containers (Phase 006):

```yaml
container:
  base_image: alpine:3.21
  registry: ghcr.io/BarunKGP/dotctl-dev
  ssh_port: 2222
  user: dev
  packages_flags: --skip-optional
```

**Fields:**
- `base_image`: Docker base image (e.g., `alpine:3.21`, `ubuntu:24.04`)
- `registry`: Container registry (e.g., `ghcr.io/user/repo`)
- `ssh_port`: SSH port to expose in container
- `user`: Non-root user to create in container
- `packages_flags`: Flags to pass to `dotctl packages install` (e.g., `--skip-optional`)

### Deploy Configuration

Configuration for remote deployments (Phase 007):

```yaml
deploy:
  k8s:
    context: homelab
    namespace: dotfiles
    storage: 10Gi
    cpu: 500m
    memory: 512Mi
```

**K8s Fields:**
- `context`: Kubernetes context name
- `namespace`: Kubernetes namespace for deployment
- `storage`: Persistent volume size
- `cpu`: CPU request (e.g., `500m`, `1`, `2`)
- `memory`: Memory request (e.g., `512Mi`, `1Gi`)

## Complete Example

See `dotctl.yaml.example` in the repository root for a complete, ready-to-use example.

## Migration from packages.conf

If you have an existing `packages.conf`, you can create a `dotctl.yaml` alongside it. Both formats are supported simultaneously. The `dotctl.yaml` takes precedence when present.

### Conversion Guide

**packages.conf:**
```
zsh       shell
neovim    editor    apk=neovim
espanso   optional  install=custom
```

**Equivalent dotctl.yaml:**
```yaml
packages:
  zsh:
    group: shell
    install:
      apt: zsh
      apk: zsh
      brew: zsh
      dnf: zsh

  neovim:
    group: editor
    install:
      apt: neovim
      apk: neovim
      brew: neovim
      dnf: neovim

  espanso:
    group: optional
    install:
      custom: scripts/packages/espanso.sh
```

## Validation

Validate your `dotctl.yaml` syntax:

```bash
dotctl doctor
```

This will check:
- Valid YAML syntax
- Required fields present (profiles, packages)
- Package references in profiles exist
- Custom install scripts are accessible

## Best Practices

1. **Keep profiles simple**: Use 3-5 profiles for common scenarios
2. **Use consistent naming**: Package names should match their canonical binary name
3. **Document custom installs**: Include comments explaining non-standard installation methods
4. **Version your config**: Commit `dotctl.yaml` to version control
5. **Use examples**: Start from `dotctl.yaml.example` and customize

## Fallback Behavior

If `dotctl.yaml` is not found:
1. dotctl continues to `packages.conf` silently
2. All functionality works as before
3. No error is raised
4. Profile flags are ignored

This ensures backward compatibility during migration.

## See Also

- `dotctl.yaml.example` - Complete example configuration
- `docs/go-cli-migration-plan.md` - Migration roadmap
- `docs/phases.md` - Implementation phases overview
