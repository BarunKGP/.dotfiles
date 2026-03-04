package config

// Config represents the dotctl.yaml configuration structure.
type Config struct {
	Profiles  map[string][]string `yaml:"profiles"`
	Packages  map[string]*Package `yaml:"packages"`
	Container *Container          `yaml:"container,omitempty"`
	Deploy    *Deploy             `yaml:"deploy,omitempty"`
}

// Package defines a package with its group and install configuration.
type Package struct {
	Group   string            `yaml:"group"`
	Install *InstallConfig    `yaml:"install"`
}

// InstallConfig specifies how a package should be installed across different package managers.
type InstallConfig struct {
	Apt          *string `yaml:"apt"`
	Apk          *string `yaml:"apk"`
	Brew         *string `yaml:"brew"`
	Dnf          *string `yaml:"dnf"`
	CustomScript *string `yaml:"custom"`  // e.g., custom: scripts/espanso.sh
}

// Container defines Docker/container image configuration (used in phase 006).
type Container struct {
	BaseImage     string `yaml:"base_image"`
	Registry      string `yaml:"registry"`
	SSHPort       int    `yaml:"ssh_port"`
	User          string `yaml:"user"`
	PackagesFlags string `yaml:"packages_flags"`
}

// Deploy defines deployment configuration (used in phase 007).
type Deploy struct {
	K8s *K8sConfig `yaml:"k8s,omitempty"`
}

// K8sConfig defines Kubernetes deployment settings.
type K8sConfig struct {
	Context   string `yaml:"context"`
	Namespace string `yaml:"namespace"`
	Storage   string `yaml:"storage"`
	CPU       string `yaml:"cpu"`
	Memory    string `yaml:"memory"`
}

// GetProfile returns the list of groups for a given profile name.
func (c *Config) GetProfile(name string) []string {
	if c == nil || c.Profiles == nil {
		return nil
	}
	return c.Profiles[name]
}

// GetPackage returns package configuration by canonical name.
func (c *Config) GetPackage(name string) *Package {
	if c == nil || c.Packages == nil {
		return nil
	}
	return c.Packages[name]
}
