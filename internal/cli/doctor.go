package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check dotfiles setup",
		Long:  "Run diagnostics to verify dotfiles installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[stub] doctor")
			return nil
		},
	}
}
