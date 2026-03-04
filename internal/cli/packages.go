package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BarunKGP/dotfiles/dotctl/internal/logging"
	"github.com/BarunKGP/dotfiles/dotctl/internal/manifest"
	"github.com/BarunKGP/dotfiles/dotctl/internal/osdetect"
	"github.com/BarunKGP/dotfiles/dotctl/internal/pkgmgr"
)

type packagesCmd struct {
	skipOptional bool
	profile      string
}

func newPackagesCmd() *cobra.Command {
	pc := &packagesCmd{}

	cmd := &cobra.Command{
		Use:   "packages",
		Short: "Manage package installation",
		Long:  "Install, list, or manage packages.",
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pc.runInstall()
		},
	}

	installCmd.Flags().BoolVar(
		&pc.skipOptional,
		"skip-optional",
		false,
		"Skip optional package groups",
	)
	installCmd.Flags().StringVar(
		&pc.profile,
		"profile",
		"default",
		"Installation profile (default, minimal, workstation, devcontainer)",
	)

	cmd.AddCommand(installCmd)

	return cmd
}

func (pc *packagesCmd) runInstall() error {
	cfg := GetGlobalConfig()
	logger := logging.NewLogger(cfg.Verbose, cfg.JSON)

	// Detect OS and package manager
	detector := osdetect.NewDetector()
	osType := detector.Detect()

	pm := pkgmgr.Detect(osType, nil)
	if pm == nil {
		return fmt.Errorf("no package manager found")
	}

	if cfg.DryRun {
		logger.DryRun(fmt.Sprintf("Package manager: %s", pm.Name()))
	} else {
		logger.Info(fmt.Sprintf("Package manager: %s", pm.Name()))
	}

	// Update package index if needed (apt)
	if pm.Name() == "apt" {
		if err := pm.UpdateIndex(cfg.DryRun); err != nil {
			return fmt.Errorf("failed to update package index: %w", err)
		}
	}

	// Parse packages.conf
	manifestPath := filepath.Join(cfg.DotfilesDir, "packages.conf")
	m, err := manifest.ParseFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to parse packages: %w", err)
	}

	// Filter packages
	filterOpts := manifest.FilterOptions{
		SkipOptional:   pc.skipOptional,
		NonInteractive: cfg.NonInteractive,
		IsInstalledFn:  pm.IsInstalled,
	}

	packages := m.Filter(filterOpts)

	if cfg.DryRun && len(packages) > 0 {
		logger.DryRun(fmt.Sprintf("Would install: %s", strings.Join(packages, ", ")))
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
			custom := pkgmgr.NewCustom(cfg.DotfilesDir, string(osType), pm.Name())
			if err := custom.Execute(pkgName, cfg.DryRun); err != nil {
				logger.Warn(fmt.Sprintf("Custom install for %s failed: %v", pkgName, err))
				continue
			}
		} else {
			// Resolve package name with overrides
			resolvedName := manifest.ResolvePackageName(*pkgInfo, pm.Name())
			if err := pm.Install(resolvedName, cfg.DryRun); err != nil {
				logger.Warn(fmt.Sprintf("Install %s failed: %v", resolvedName, err))
			}
		}
	}

	logger.Done("Package installation complete")
	return nil
}
