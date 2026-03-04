package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/BarunKGP/dotfiles/dotctl/internal/logging"
	"github.com/BarunKGP/dotfiles/dotctl/internal/stow"
)

type stowCmd struct {
	minimal bool
}

func newStowCmd() *cobra.Command {
	sc := &stowCmd{}

	cmd := &cobra.Command{
		Use:   "stow",
		Short: "Manage stow operations",
		Long:  "Apply or check stow operations.",
	}

	applyCmd := &cobra.Command{
		Use:   "apply [packages...]",
		Short: "Apply stow links",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.runApply(args)
		},
	}

	applyCmd.Flags().BoolVar(
		&sc.minimal,
		"minimal",
		false,
		"Minimal stow (skip nvim, tmux)",
	)

	cmd.AddCommand(applyCmd)

	return cmd
}

func (sc *stowCmd) runApply(args []string) error {
	cfg := GetGlobalConfig()
	logger := logging.NewLogger(cfg.Verbose, cfg.JSON)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	stower := stow.New(cfg.DotfilesDir, home)

	// Use provided packages or defaults
	var pkgs []string
	if len(args) > 0 {
		pkgs = args
	} else {
		pkgs = stow.GetDefaultPackages(sc.minimal)
	}

	if cfg.DryRun && len(pkgs) > 0 {
		logger.DryRun(fmt.Sprintf("Would stow: %s", strings.Join(pkgs, ", ")))
	}

	for _, pkg := range pkgs {
		if err := stower.Apply(pkg, cfg.DryRun); err != nil {
			logger.Fail(fmt.Sprintf("stow %s failed: %v", pkg, err))
			return err
		}
	}

	logger.Done("Stow complete")
	return nil
}
