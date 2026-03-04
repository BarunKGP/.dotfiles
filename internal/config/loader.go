package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load attempts to load dotctl.yaml from the given directory.
// Returns nil if file not found (fallback to packages.conf).
// Returns error if file exists but cannot be parsed.
func Load(dotfilesDir string) (*Config, error) {
	configPath := filepath.Join(dotfilesDir, "dotctl.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil  // Fallback to packages.conf
		}
		return nil, fmt.Errorf("failed to read dotctl.yaml: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse dotctl.yaml: %w", err)
	}

	return &cfg, nil
}

// Validate performs basic validation on the config.
func (c *Config) Validate() error {
	if c == nil {
		return nil  // No config file is valid (fallback case)
	}

	if c.Profiles == nil || len(c.Profiles) == 0 {
		return fmt.Errorf("no profiles defined in dotctl.yaml")
	}

	if c.Packages == nil || len(c.Packages) == 0 {
		return fmt.Errorf("no packages defined in dotctl.yaml")
	}

	return nil
}
