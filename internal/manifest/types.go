package manifest

// Package represents a package entry in packages.conf.
type Package struct {
	Canonical    string            // canonical name (e.g., "zsh", "neovim")
	Group        string            // group (e.g., "shell", "editor", "optional")
	PMOverrides   map[string]string // pm-specific overrides (e.g., apk=neovim)
	CustomInstall bool              // whether to use custom install script
}

// Manifest is the parsed packages.conf.
type Manifest struct {
	Packages []Package
}
