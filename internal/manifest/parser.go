package manifest

import (
	"fmt"
	"os"
	"strings"
)

// ParseFile parses a packages.conf file.
func ParseFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	return Parse(string(data))
}

// Parse parses the content of packages.conf.
func Parse(content string) (*Manifest, error) {
	manifest := &Manifest{
		Packages: []Package{},
	}

	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		pkg, err := parseLine(line, lineNum)
		if err != nil {
			return nil, err
		}

		if pkg != nil {
			manifest.Packages = append(manifest.Packages, *pkg)
		}
	}

	return manifest, nil
}

// parseLine parses a single line from packages.conf.
// Format: canonical group [key=value...]
// Example: zsh shell
// Example: neovim editor apk=neovim
// Example: espanso optional install=custom
func parseLine(line string, lineNum int) (*Package, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil, fmt.Errorf("line %d: expected at least 2 fields (canonical, group)", lineNum+1)
	}

	canonical := fields[0]
	group := fields[1]

	pkg := &Package{
		Canonical:    canonical,
		Group:        group,
		PMOverrides:   make(map[string]string),
	}

	// Parse optional key=value tokens
	for _, token := range fields[2:] {
		if strings.Contains(token, "=") {
			parts := strings.SplitN(token, "=", 2)
			key := parts[0]
			value := parts[1]

			switch key {
			case "install":
				if value == "custom" {
					pkg.CustomInstall = true
				}
			// Handle package manager overrides (e.g., apk=neovim)
			case "apt", "apk", "brew", "dnf", "pacman":
				pkg.PMOverrides[key] = value
			}
		}
	}

	return pkg, nil
}
