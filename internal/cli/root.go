package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/BarunKGP/dotfiles/dotctl/internal/config"
)

type GlobalConfig struct {
	DotfilesDir    string
	NonInteractive bool
	DryRun         bool
	Verbose        bool
	JSON           bool
	Config         *config.Config
}

var globalConfig GlobalConfig

// NewRootCmd creates the root command for dotctl.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dotctl",
		Short: "dotctl - dotfiles installer and manager",
		Long:  "A Go-based CLI for managing dotfiles installation, stowing, and environment setup.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return resolveDotfilesDir()
		},
	}

	// Global flags
	cmd.PersistentFlags().StringVar(
		&globalConfig.DotfilesDir,
		"dotfiles-dir",
		"",
		"Path to dotfiles directory (auto-detected if not specified)",
	)
	cmd.PersistentFlags().BoolVar(
		&globalConfig.NonInteractive,
		"non-interactive",
		false,
		"Skip interactive prompts",
	)
	cmd.PersistentFlags().BoolVar(
		&globalConfig.DryRun,
		"dry-run",
		false,
		"Show what would be done without making changes",
	)
	cmd.PersistentFlags().BoolVar(
		&globalConfig.Verbose,
		"verbose",
		false,
		"Enable verbose output",
	)
	cmd.PersistentFlags().BoolVar(
		&globalConfig.JSON,
		"json",
		false,
		"Output in JSON format",
	)

	// Add subcommands
	cmd.AddCommand(newInstallCmd())
	cmd.AddCommand(newPackagesCmd())
	cmd.AddCommand(newStowCmd())
	cmd.AddCommand(newEnvCmd())
	cmd.AddCommand(newDoctorCmd())

	return cmd
}

// resolveDotfilesDir resolves the dotfiles directory by:
// 1. Using the --dotfiles-dir flag if provided
// 2. Walking up from the executable's directory looking for packages.conf
// 3. Falling back to the current working directory
func resolveDotfilesDir() error {
	if globalConfig.DotfilesDir != "" {
		// Flag was provided, use it as-is
		return nil
	}

	// Try to find packages.conf by walking up from executable
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		if found := findDotfilesDir(exeDir); found != "" {
			globalConfig.DotfilesDir = found
			return nil
		}
	}

	// Try walking up from cwd
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	if found := findDotfilesDir(cwd); found != "" {
		globalConfig.DotfilesDir = found
		return nil
	}

	// Fall back to cwd
	globalConfig.DotfilesDir = cwd

	// Try to load dotctl.yaml (non-fatal if not found)
	cfg, err := config.Load(globalConfig.DotfilesDir)
	if err != nil {
		// Log warning but don't fail - packages.conf is fallback
		if globalConfig.Verbose {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load dotctl.yaml: %v\n", err)
		}
	}
	globalConfig.Config = cfg

	return nil
}

// findDotfilesDir walks up from dir looking for packages.conf.
func findDotfilesDir(dir string) string {
	current := dir
	for {
		configPath := filepath.Join(current, "packages.conf")
		if _, err := os.Stat(configPath); err == nil {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached root
			break
		}
		current = parent
	}
	return ""
}

// Execute runs the CLI.
func Execute() error {
	cmd := NewRootCmd()
	return cmd.Execute()
}

// GetGlobalConfig returns the resolved global configuration.
func GetGlobalConfig() GlobalConfig {
	return globalConfig
}

// HasSudo checks if the process is running as root or if sudo is available.
func HasSudo() (bool, error) {
	if os.Getuid() == 0 {
		return true, nil
	}
	_, err := exec.LookPath("sudo")
	return err == nil, nil
}
