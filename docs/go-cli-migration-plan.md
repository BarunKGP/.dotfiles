# Dotfiles Installer Migration Plan: Bash Scripts → Go CLI

## Why migrate

The current bootstrap flow is split across `install.sh`, `scripts/install-packages.sh`, package-specific scripts, and shell modules. This works, but orchestration logic (flags, OS detection, package manager branching, optional prompts, stow operations, shell defaults, local config generation) is increasingly complex to evolve safely. A Go CLI provides:

- stronger structure for complex workflows,
- testable business logic without shell side effects,
- clearer extensibility for future language tooling support,
- better UX consistency (subcommands, validation, dry-run, structured output).

## Current behavior to preserve

The Go CLI should fully preserve functional behavior before adding new features:

1. Parse installer flags equivalent to:
   - `--minimal`
   - `--no-packages`
   - `--skip-optional`
2. Detect OS as one of `wsl2`, `macos`, `linux`.
3. Detect package manager (`brew`, `apt`, `apk`, `dnf`, unknown).
4. Parse and apply `packages.conf` semantics:
   - group handling (`util`, `editor`, `shell`, `terminal`, `optional`),
   - per-manager package rename tokens (e.g. `apt=...`),
   - `install=custom` package hooks.
5. Support non-interactive environments by auto-skipping optional packages.
6. Run stow for core packages and minimal/full variants.
7. Attempt default shell switch to zsh when available.
8. Create `~/.config/shell/local.sh` template if missing.
9. Print OS-specific post-install guidance.

## Target CLI design

### Binary and command model

Recommended binary name: `dotctl` (or `dotfiles` if preferred).

Initial command tree:

- `dotctl install` → full machine bootstrap (replacement for `./install.sh`)
- `dotctl packages install` → package install only
- `dotctl stow apply` → stow symlinks only
- `dotctl env init` → create/merge local env config template
- `dotctl doctor` → verify required dependencies and config health

Common global flags:

- `--dotfiles-dir` (default: executable-relative repo root or cwd)
- `--non-interactive`
- `--dry-run`
- `--verbose`
- `--json` (optional, for machine-readable diagnostics)

### Internal package layout (Go)

```text
cmd/dotctl/
  main.go
internal/
  cli/            # cobra/urfave command wiring
  install/        # top-level install orchestration
  osdetect/       # wsl2/macos/linux detection
  pkgmgr/         # package manager abstraction + implementations
  manifest/       # packages.conf parser/validator
  stow/           # stow operation wrapper
  shell/          # shell detection + default shell changes
  envconfig/      # local.sh generation + merge behavior
  languages/      # language runtime definitions (future-facing)
  logging/        # console formatting + structured logs
```

## Config evolution plan

### 1) Keep `packages.conf` initially (compatibility phase)

Start by reusing existing `packages.conf` and parsing it in Go to avoid a risky format migration during bootstrap replacement.

### 2) Introduce `dotctl.yaml` for richer behavior (phase 2)

Add a typed config format that can represent:

- packages by capability and profile,
- language runtime setup,
- environment variable templates,
- OS/package-manager overrides,
- optional interactive prompts.

Example high-level schema:

```yaml
profiles:
  default: [core, shell, editor]
  minimal: [core]

packages:
  - name: zsh
    group: shell
    managers:
      apt: zsh
      brew: zsh

languages:
  go:
    strategy: system
    env:
      GOROOT: /usr/local/go
      GOPATH: "$HOME/go"
      add_path: ["$GOPATH/bin", "$GOROOT/bin"]
  python:
    strategy: pyenv
    env:
      PYENV_ROOT: "$HOME/.pyenv"
      add_path: ["$PYENV_ROOT/bin"]
```

## Language support strategy (future requirement)

To support languages such as Go/Python/Node in a maintainable way, add a provider model:

- `languages.Provider` interface:
  - `Detect(ctx)`, `Install(ctx)`, `Env(ctx)`, `Validate(ctx)`
- provider implementations:
  - `go_system` (GOROOT/GOPATH conventions),
  - `python_pyenv` (PYENV_ROOT + pyenv init guidance),
  - `node_nvm` (NVM_DIR and shell integration).

### Environment variable automation

Use a generated file model instead of hand-editing shell modules:

- Generated output path: `~/.config/shell/generated/languages.sh`
- `dotctl env render` regenerates this file from declarative config.
- Existing shell entry points source generated files if present.

This allows controlled updates for variables like:

- `GOROOT`, `GOPATH`
- `PYENV_ROOT` / `PYENV_HOME` (project preference)
- `NVM_DIR`
- language-specific PATH segments

and avoids duplicated logic across bash/zsh modules.

## Migration phases

### Phase 0 — repo preparation

- Add Go module scaffold (`go.mod`, `cmd/dotctl/main.go`).
- Add CI job for `go test` and `go vet`.
- Add `dotctl install --dry-run` skeleton with parity check output.

### Phase 1 — behavior parity with existing install flow

- Implement command flags and OS/package-manager detection.
- Implement `packages.conf` parser and installer execution.
- Implement stow orchestration and local config template generation.
- Keep existing shell scripts as fallback until parity is verified.

Exit criteria:

- `dotctl install` matches `install.sh` behavior for minimal/full/no-packages modes.

### Phase 2 — language automation and generated env

- Add `languages` provider framework.
- Add `dotctl env render` to generate language environment file(s).
- Update `shell/common/env.sh` and `shell/common/path.sh` to source generated files.

Exit criteria:

- Go and Python env variables/path can be declared once and applied automatically.

### Phase 3 — config modernization

- Introduce `dotctl.yaml` with migration tooling from `packages.conf`.
- Support profiles (default/minimal/workstation/devcontainer).
- Add richer package metadata and optional prerequisites.

### Phase 4 — deprecate bash installer

- Convert `install.sh` into thin shim that calls `dotctl install`.
- Mark shell scripts as legacy and remove duplicated orchestration logic.

## Testing and validation plan

### Unit tests

- OS detection matrix tests.
- Package manager selection tests.
- `packages.conf` parsing/validation tests.
- Group/flag decision tests (`minimal`, `optional`, non-interactive).
- Env rendering golden tests for generated shell output.

### Integration tests

- Dry-run integration snapshots for linux/macos/wsl2 detection stubs.
- Containerized smoke tests for apt/apk flows.
- Idempotency test: running install twice should converge without failures.

### Safety checks

- `dotctl doctor` verifies required commands (`stow`, shell tools).
- rollback-friendly behavior with explicit logging around each step.

## Immediate concrete changes recommended next

1. Create Go scaffold and `dotctl install --dry-run` command.
2. Port OS + package manager detection from bash into tested Go utilities.
3. Port `packages.conf` parser and decision logic (including optional handling).
4. Port stow and local config template generation.
5. Add generated env file path and wire shell modules to source it when present.

This sequence minimizes risk while unblocking the future language-support roadmap.
