package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPackagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages",
		Short: "Manage package installation",
		Long:  "Install, list, or manage packages.",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "install",
		Short: "Install packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[stub] packages install")
			return nil
		},
	})

	return cmd
}
