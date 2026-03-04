package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/BarunKGP/dotfiles/dotctl/internal/envconfig"
	"github.com/BarunKGP/dotfiles/dotctl/internal/languages"
	"github.com/BarunKGP/dotfiles/dotctl/internal/logging"
)

func newEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Manage environment configuration",
		Long:  "Initialize or configure environment variables.",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnvInit()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "render",
		Short: "Generate language environment variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnvRender()
		},
	})

	return cmd
}

func runEnvInit() error {
	cfg := GetGlobalConfig()
	logger := logging.NewLogger(cfg.Verbose, cfg.JSON)

	// Resolve XDG_CONFIG_HOME
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		xdgHome = filepath.Join(home, ".config")
	}

	created, configPath, err := envconfig.EnsureLocalConfig(xdgHome, cfg.DryRun)
	if err != nil {
		return err
	}

	if cfg.DryRun {
		logger.DryRun(fmt.Sprintf("Would create: %s", configPath))
	} else if created {
		logger.Done(fmt.Sprintf("Created local config: %s", configPath))
	} else {
		logger.Skip(fmt.Sprintf("Local config already exists: %s", configPath))
	}

	return nil
}

func runEnvRender() error {
	cfg := GetGlobalConfig()
	logger := logging.NewLogger(cfg.Verbose, cfg.JSON)

	// Resolve XDG_CONFIG_HOME
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		xdgHome = filepath.Join(home, ".config")
	}

	// Generate language environment
	registry := languages.NewRegistry()
	script, err := registry.Render(context.Background())
	if err != nil {
		return fmt.Errorf("failed to render language environment: %w", err)
	}

	if cfg.DryRun {
		fmt.Println(script)
		logger.DryRun("Would generate language environment")
		return nil
	}

	// Create generated directory
	generatedDir := filepath.Join(xdgHome, "shell", "generated")
	if err := os.MkdirAll(generatedDir, 0755); err != nil {
		return fmt.Errorf("failed to create generated directory: %w", err)
	}

	// Write generated file
	outputPath := filepath.Join(generatedDir, "languages.sh")
	if err := os.WriteFile(outputPath, []byte(script), 0644); err != nil {
		return fmt.Errorf("failed to write generated file: %w", err)
	}

	logger.Done(fmt.Sprintf("Generated language environment: %s", outputPath))
	return nil
}
