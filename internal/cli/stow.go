package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stow",
		Short: "Manage stow operations",
		Long:  "Apply or check stow operations.",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Apply stow links",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[stub] stow apply")
			return nil
		},
	})

	return cmd
}
