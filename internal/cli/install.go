package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/BarunKGP/dotfiles/dotctl/internal/install"
	"github.com/BarunKGP/dotfiles/dotctl/internal/logging"
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

	logger := logging.NewLogger(cfg.Verbose, cfg.JSON)

	opts := install.Options{
		Minimal:        ic.minimal,
		NoPackages:     ic.noPackages,
		SkipOptional:   ic.skipOptional,
		DotfilesDir:    cfg.DotfilesDir,
		NonInteractive: cfg.NonInteractive,
		DryRun:         cfg.DryRun,
		Verbose:        cfg.Verbose,
	}

	installer := install.New(opts, logger)
	return installer.Run(ctx)
}
