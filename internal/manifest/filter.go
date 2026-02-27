package manifest

// DecisionType represents what should be done with a package.
type DecisionType int

const (
	DecisionInstall DecisionType = iota
	DecisionSkip
	DecisionPrompt
	DecisionAlreadyInstalled
)

// FilterOptions controls package filtering.
type FilterOptions struct {
	Minimal        bool
	SkipOptional   bool
	NonInteractive bool
	IsInstalledFn  func(pkg string) bool
}

// Filter returns packages to install based on options.
func (m *Manifest) Filter(opts FilterOptions) []string {
	var result []string

	for _, pkg := range m.Packages {
		// Decide on this package
		decision := DecidePackage(pkg, opts)

		if decision == DecisionInstall {
			result = append(result, pkg.Canonical)
		}
	}

	return result
}

// DecidePackage determines what to do with a single package.
func DecidePackage(pkg Package, opts FilterOptions) DecisionType {
	// Skip if already installed and not forced
	if opts.IsInstalledFn != nil && opts.IsInstalledFn(pkg.Canonical) {
		return DecisionAlreadyInstalled
	}

	// Skip optional packages if requested
	if opts.SkipOptional && pkg.Group == "optional" {
		return DecisionSkip
	}

	// Skip certain groups in minimal mode
	if opts.Minimal && (pkg.Group == "editor" || pkg.Group == "terminal") {
		return DecisionSkip
	}

	// In non-interactive mode, skip prompts (install all others)
	if opts.NonInteractive {
		return DecisionInstall
	}

	// Interactive: prompt for uncertain cases
	if pkg.Group == "optional" {
		return DecisionPrompt
	}

	return DecisionInstall
}

// ResolvePackageName returns the name to use for a package given a package manager.
func ResolvePackageName(pkg Package, pmName string) string {
	// Check for package manager override
	if override, ok := pkg.PMOverrides[pmName]; ok {
		return override
	}

	// Fall back to canonical name
	return pkg.Canonical
}
