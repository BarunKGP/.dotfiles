package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type installCmd struct {
	minimal        bool
	noPackages     bool
	skipOptional   bool
}

func newInstallCmd() *cobra.Command {
	ic := &installCmd{}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install and configure dotfiles",
		Long:  "Install packages, stow dotfiles, and set up shell configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ic.run(cmd.Context())
		},
	}

	cmd.Flags().BoolVar(
		&ic.minimal,
		"minimal",
		false,
		"Minimal install (skip nvim, tmux)",
	)
	cmd.Flags().BoolVar(
		&ic.noPackages,
		"no-packages",
		false,
		"Skip package installation",
	)
	cmd.Flags().BoolVar(
		&ic.skipOptional,
		"skip-optional",
		false,
		"Skip optional package groups",
	)

	return cmd
}

func (ic *installCmd) run(ctx context.Context) error {
	cfg := GetGlobalConfig()

	// Detect OS
	osType := detectOS()
	if cfg.DryRun {
		fmt.Printf("[dry-run] OS detection: %s\n", osType)
	}

	// Detect package manager
	pm := detectPackageManager(osType)
	if cfg.DryRun {
		fmt.Printf("[dry-run] Package manager: %s\n", pm)
	}

	// Determine packages to install
	if !ic.noPackages {
		packages := resolvePackages(cfg.DotfilesDir, ic.minimal, ic.skipOptional)
		if cfg.DryRun && len(packages) > 0 {
			fmt.Printf("[dry-run] Would install: %s\n", strings.Join(packages, ", "))
		}
	}

	// Determine stow packages
	stowPackages := resolveStowPackages(ic.minimal)
	if cfg.DryRun && len(stowPackages) > 0 {
		fmt.Printf("[dry-run] Would stow: %s\n", strings.Join(stowPackages, ", "))
	}

	// Local config
	if cfg.DryRun {
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".config", "shell", "local.sh")
		fmt.Printf("[dry-run] Would create: %s\n", configPath)
	}

	if cfg.DryRun {
		return nil
	}

	// Non-dry-run logic would go here (Phase 1)
	return nil
}

// detectOS returns the detected OS type (linux, wsl2, macos).
func detectOS() string {
	// Check for WSL2
	if data, err := os.ReadFile("/proc/version"); err == nil {
		if strings.Contains(string(data), "microsoft") || strings.Contains(string(data), "WSL") {
			return "wsl2"
		}
	}

	// Check for macOS
	if _, err := os.Stat("/System/Library/CoreServices/Finder.app"); err == nil {
		return "macos"
	}

	// Default to linux
	return "linux"
}

// detectPackageManager returns the detected package manager (brew, apt, apk, dnf).
func detectPackageManager(osType string) string {
	managers := []string{"brew", "apt", "apk", "dnf", "pacman"}

	if osType == "macos" {
		managers = []string{"brew"}
	} else if osType == "wsl2" {
		managers = []string{"apt", "brew"}
	}

	for _, mgr := range managers {
		if pathExists(mgr) {
			return mgr
		}
	}

	return "unknown"
}

// pathExists checks if a command exists in PATH.
func pathExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// resolvePackages returns the list of packages to install.
func resolvePackages(dotfilesDir string, minimal, skipOptional bool) []string {
	// Parse packages.conf
	confPath := filepath.Join(dotfilesDir, "packages.conf")
	data, err := os.ReadFile(confPath)
	if err != nil {
		return nil
	}

	var packages []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		canonical := fields[0]
		group := fields[1]

		// Skip based on flags
		if skipOptional && group == "optional" {
			continue
		}

		// Skip shell, editor, terminal if minimal
		if minimal && (group == "shell" || group == "editor" || group == "terminal") {
			continue
		}

		packages = append(packages, canonical)
	}

	return packages
}

// resolveStowPackages returns the list of packages to stow.
func resolveStowPackages(minimal bool) []string {
	packages := []string{"bash", "git", "nvim"}

	if !minimal {
		packages = append(packages, "zsh", "tmux")
	}

	return packages
}
