package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/BarunKGP/dotfiles/dotctl/internal/doctor"
	"github.com/BarunKGP/dotfiles/dotctl/internal/logging"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check dotfiles setup",
		Long:  "Run diagnostics to verify dotfiles installation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor()
		},
	}
}

func runDoctor() error {
	cfg := GetGlobalConfig()
	logger := logging.NewLogger(cfg.Verbose, cfg.JSON)

	reports := doctor.Diagnose(cfg.DotfilesDir)

	hasFailure := false
	for _, report := range reports {
		switch report.Status {
		case "PASS":
			logger.Done(fmt.Sprintf("%s: %s", report.Name, report.Details))
		case "FAIL":
			logger.Fail(fmt.Sprintf("%s: %s", report.Name, report.Details))
			hasFailure = true
		case "WARN":
			logger.Warn(fmt.Sprintf("%s: %s", report.Name, report.Details))
		}
	}

	if hasFailure {
		return fmt.Errorf("diagnostic failures detected")
	}

	return nil
}
