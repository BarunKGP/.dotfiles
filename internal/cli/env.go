package cli

import (
	"fmt"

	"github.com/spf13/cobra"
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
			fmt.Println("[stub] env init")
			return nil
		},
	})

	return cmd
}
