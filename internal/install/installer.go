package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BarunKGP/dotfiles/dotctl/internal/envconfig"
	"github.com/BarunKGP/dotfiles/dotctl/internal/languages"
	"github.com/BarunKGP/dotfiles/dotctl/internal/logging"
	"github.com/BarunKGP/dotfiles/dotctl/internal/manifest"
	"github.com/BarunKGP/dotfiles/dotctl/internal/osdetect"
	"github.com/BarunKGP/dotfiles/dotctl/internal/pkgmgr"
	"github.com/BarunKGP/dotfiles/dotctl/internal/shell"
	"github.com/BarunKGP/dotfiles/dotctl/internal/stow"
)

// Installer orchestrates the installation process.
type Installer struct {
	Opts   Options
	Logger *logging.Logger
	OSType osdetect.OSType
	PM     pkgmgr.Manager
}

// New creates a new Installer.
func New(opts Options, logger *logging.Logger) *Installer {
	return &Installer{
		Opts:   opts,
		Logger: logger,
	}
}

// Run executes the installation.
func (i *Installer) Run(ctx context.Context) error {
	// Step 1: Print header
	i.printHeader()

	// Step 2: Install packages if not --no-packages
	if !i.Opts.NoPackages {
		if err := i.installPackages(); err != nil {
			i.Logger.Fail(fmt.Sprintf("package installation failed: %v", err))
			if !i.Opts.DryRun {
				return err
			}
		}
	}

	// Step 3: Check stow available
	if err := stow.Check(); err != nil {
		i.Logger.Fail(fmt.Sprintf("stow check failed: %v", err))
		if !i.Opts.DryRun {
			return err
		}
	}

	// Step 4: Stow packages
	if err := i.stowPackages(); err != nil {
		i.Logger.Fail(fmt.Sprintf("stow failed: %v", err))
		if !i.Opts.DryRun {
			return err
		}
	}

	// Step 5: Set default shell to zsh if available
	availableShells := shell.AvailableShells(nil)
	hasZsh := false
	for _, s := range availableShells {
		if s == "zsh" {
			hasZsh = true
			break
		}
	}

	if hasZsh && !i.Opts.DryRun {
		// Non-fatal if chsh fails
		result := shell.SetDefault("zsh")
		if result.Success {
			i.Logger.Done("Set default shell to zsh")
		} else {
			i.Logger.Warn(fmt.Sprintf("Could not set zsh as default: %v", result.Error))
		}
	}

	// Step 6: Ensure local config
	if err := i.ensureLocalConfig(); err != nil {
		i.Logger.Fail(fmt.Sprintf("local config creation failed: %v", err))
		if !i.Opts.DryRun {
			return err
		}
	}

	// Step 7: Render language environment
	if err := i.renderLanguageEnv(); err != nil {
		i.Logger.Warn(fmt.Sprintf("language env render failed: %v", err))
		// non-fatal — user can run: dotctl env render
	}

	// Step 8: Print post-install guidance
	i.printPostInstallGuidance()

	// Step 9: Print next steps
	i.printNextSteps()

	return nil
}

func (i *Installer) printHeader() {
	detector := osdetect.NewDetector()
	i.OSType = detector.Detect()

	mode := "install"
	if i.Opts.DryRun {
		mode = "dry-run"
	}

	i.Logger.Section("Installation")
	if i.Opts.DryRun {
		i.Logger.DryRun(fmt.Sprintf("OS detection: %s", i.OSType))
	} else {
		i.Logger.Info(fmt.Sprintf("OS detection: %s", i.OSType))
	}
	i.Logger.Info(fmt.Sprintf("Mode: %s", mode))
	i.Logger.Info(fmt.Sprintf("Dotfiles dir: %s", i.Opts.DotfilesDir))
}

func (i *Installer) installPackages() error {
	// Detect package manager
	i.PM = pkgmgr.Detect(i.OSType, nil)
	if i.PM == nil {
		return fmt.Errorf("no package manager found")
	}

	if i.Opts.DryRun {
		i.Logger.DryRun(fmt.Sprintf("Package manager: %s", i.PM.Name()))
	} else {
		i.Logger.Info(fmt.Sprintf("Package manager: %s", i.PM.Name()))
	}

	// Update package index if needed (apt)
	if i.PM.Name() == "apt" {
		if err := i.PM.UpdateIndex(i.Opts.DryRun); err != nil {
			return fmt.Errorf("failed to update package index: %w", err)
		}
	}

	// Parse packages.conf
	manifestPath := filepath.Join(i.Opts.DotfilesDir, "packages.conf")
	m, err := manifest.ParseFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to parse packages: %w", err)
	}

	// Filter packages
	filterOpts := manifest.FilterOptions{
		Minimal:        i.Opts.Minimal,
		SkipOptional:   i.Opts.SkipOptional,
		NonInteractive: i.Opts.NonInteractive,
		IsInstalledFn:  i.PM.IsInstalled,
	}

	packages := m.Filter(filterOpts)

	if i.Opts.DryRun && len(packages) > 0 {
		i.Logger.DryRun(fmt.Sprintf("Would install: %s", strings.Join(packages, ", ")))
	}

	// Install each package
	for _, pkgName := range packages {
		// Find the package in the manifest
		var pkgInfo *manifest.Package
		for _, p := range m.Packages {
			if p.Canonical == pkgName {
				pkgInfo = &p
				break
			}
		}

		if pkgInfo == nil {
			continue
		}

		// Check if custom install
		if pkgInfo.CustomInstall {
			custom := pkgmgr.NewCustom(i.Opts.DotfilesDir, string(i.OSType), i.PM.Name())
			if err := custom.Execute(pkgName, i.Opts.DryRun); err != nil {
				i.Logger.Warn(fmt.Sprintf("Custom install for %s failed: %v", pkgName, err))
				continue
			}
		} else {
			// Resolve package name with overrides
			resolvedName := manifest.ResolvePackageName(*pkgInfo, i.PM.Name())
			if err := i.PM.Install(resolvedName, i.Opts.DryRun); err != nil {
				i.Logger.Warn(fmt.Sprintf("Install %s failed: %v", resolvedName, err))
			}
		}
	}

	return nil
}

func (i *Installer) stowPackages() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	stower := stow.New(i.Opts.DotfilesDir, home)
	pkgs := stow.GetDefaultPackages(i.Opts.Minimal)

	if i.Opts.DryRun && len(pkgs) > 0 {
		i.Logger.DryRun(fmt.Sprintf("Would stow: %s", strings.Join(pkgs, ", ")))
	}

	for _, pkg := range pkgs {
		if err := stower.Apply(pkg, i.Opts.DryRun); err != nil {
			return fmt.Errorf("stow %s failed: %w", pkg, err)
		}
	}

	return nil
}

func (i *Installer) ensureLocalConfig() error {
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		xdgHome = filepath.Join(home, ".config")
	}

	created, configPath, err := envconfig.EnsureLocalConfig(xdgHome, i.Opts.DryRun)
	if err != nil {
		return err
	}

	if i.Opts.DryRun {
		i.Logger.DryRun(fmt.Sprintf("Would create: %s", configPath))
	} else if created {
		i.Logger.Done(fmt.Sprintf("Created local config: %s", configPath))
	} else {
		i.Logger.Skip(fmt.Sprintf("Local config already exists: %s", configPath))
	}

	return nil
}

func (i *Installer) renderLanguageEnv() error {
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		xdgHome = filepath.Join(home, ".config")
	}

	registry := languages.NewRegistry()
	script, err := registry.Render(context.Background())
	if err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	if i.Opts.DryRun {
		i.Logger.DryRun("Would render language environment")
		return nil
	}

	generatedDir := filepath.Join(xdgHome, "shell", "generated")
	if err := os.MkdirAll(generatedDir, 0755); err != nil {
		return err
	}
	outputPath := filepath.Join(generatedDir, "languages.sh")
	if err := os.WriteFile(outputPath, []byte(script), 0644); err != nil {
		return err
	}
	i.Logger.Done(fmt.Sprintf("Generated language environment: %s", outputPath))
	return nil
}

func (i *Installer) printPostInstallGuidance() {
	i.Logger.Section("Post-Install Guidance")

	switch i.OSType {
	case osdetect.WSL2:
		i.Logger.Info("WSL2: Espanso will be installed on the Windows host")
		i.Logger.Info("See scripts/packages/espanso.sh for details")
	case osdetect.MacOS:
		i.Logger.Info("macOS: Homebrew installed packages via brew")
	case osdetect.Linux:
		i.Logger.Info("Linux: Packages installed via system package manager")
	}
}

func (i *Installer) printNextSteps() {
	i.Logger.Section("Next Steps")
	i.Logger.Info("1. Reload your shell: exec $SHELL")
	i.Logger.Info("2. Check dotfiles setup: dotctl doctor")
	i.Logger.Info("3. Customize local config: $XDG_CONFIG_HOME/shell/local.sh")
}
