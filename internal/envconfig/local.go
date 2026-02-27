package envconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureLocalConfig creates the local shell config file if it doesn't exist.
// Returns (created, path, error).
func EnsureLocalConfig(xdgHome string, dryRun bool) (bool, string, error) {
	if xdgHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return false, "", fmt.Errorf("failed to get home directory: %w", err)
		}
		xdgHome = filepath.Join(home, ".config")
	}

	configDir := filepath.Join(xdgHome, "shell")
	configPath := filepath.Join(configDir, "local.sh")

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		// File exists, idempotent
		return false, configPath, nil
	}

	if dryRun {
		return true, configPath, nil
	}

	// Create directory if needed
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return false, "", fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create the local.sh file with a template
	template := `#!/bin/bash
# Local shell configuration
# This file is gitignored and machine-specific

# Add your local environment variables, secrets, and custom settings here
# Example:
# export CUSTOM_VAR="value"
`

	if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
		return false, "", fmt.Errorf("failed to create local.sh: %w", err)
	}

	return true, configPath, nil
}
