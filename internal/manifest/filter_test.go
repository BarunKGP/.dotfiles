package manifest

import (
	"testing"
)

func TestDecidePackage(t *testing.T) {
	tests := []struct {
		name     string
		pkg      Package
		opts     FilterOptions
		expected DecisionType
	}{
		{
			name: "normal package",
			pkg: Package{
				Canonical: "bash",
				Group:     "shell",
			},
			opts:     FilterOptions{},
			expected: DecisionInstall,
		},
		{
			name: "already installed",
			pkg: Package{
				Canonical: "bash",
				Group:     "shell",
			},
			opts: FilterOptions{
				IsInstalledFn: func(s string) bool { return s == "bash" },
			},
			expected: DecisionAlreadyInstalled,
		},
		{
			name: "optional skipped",
			pkg: Package{
				Canonical: "espanso",
				Group:     "optional",
			},
			opts: FilterOptions{
				SkipOptional: true,
			},
			expected: DecisionSkip,
		},
		{
			name: "editor in minimal mode",
			pkg: Package{
				Canonical: "neovim",
				Group:     "editor",
			},
			opts: FilterOptions{
				Minimal: true,
			},
			expected: DecisionSkip,
		},
		{
			name: "terminal in minimal mode",
			pkg: Package{
				Canonical: "tmux",
				Group:     "terminal",
			},
			opts: FilterOptions{
				Minimal: true,
			},
			expected: DecisionSkip,
		},
		{
			name: "optional in non-interactive",
			pkg: Package{
				Canonical: "espanso",
				Group:     "optional",
			},
			opts: FilterOptions{
				NonInteractive: true,
			},
			expected: DecisionInstall,
		},
		{
			name: "optional in interactive",
			pkg: Package{
				Canonical: "espanso",
				Group:     "optional",
			},
			opts:     FilterOptions{},
			expected: DecisionPrompt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecidePackage(tt.pkg, tt.opts)
			if got != tt.expected {
				t.Errorf("DecidePackage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResolvePackageName(t *testing.T) {
	tests := []struct {
		name     string
		pkg      Package
		pmName   string
		expected string
	}{
		{
			name: "no override",
			pkg: Package{
				Canonical: "neovim",
			},
			pmName:   "apt",
			expected: "neovim",
		},
		{
			name: "with apk override",
			pkg: Package{
				Canonical:  "neovim",
				PMOverrides: map[string]string{"apk": "neovim"},
			},
			pmName:   "apk",
			expected: "neovim",
		},
		{
			name: "different pm no override",
			pkg: Package{
				Canonical:  "neovim",
				PMOverrides: map[string]string{"apk": "neovim"},
			},
			pmName:   "apt",
			expected: "neovim",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolvePackageName(tt.pkg, tt.pmName)
			if got != tt.expected {
				t.Errorf("ResolvePackageName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestManifestFilter(t *testing.T) {
	manifest := &Manifest{
		Packages: []Package{
			{Canonical: "bash", Group: "shell"},
			{Canonical: "zsh", Group: "shell"},
			{Canonical: "neovim", Group: "editor"},
			{Canonical: "espanso", Group: "optional"},
		},
	}

	tests := []struct {
		name     string
		opts     FilterOptions
		expected []string
	}{
		{
			name: "all packages",
			opts: FilterOptions{
				NonInteractive: true,
			},
			expected: []string{"bash", "zsh", "neovim", "espanso"},
		},
		{
			name: "skip optional",
			opts: FilterOptions{
				SkipOptional:   true,
				NonInteractive: true,
			},
			expected: []string{"bash", "zsh", "neovim"},
		},
		{
			name: "minimal",
			opts: FilterOptions{
				Minimal:        true,
				NonInteractive: true,
			},
			expected: []string{"bash", "zsh", "espanso"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manifest.Filter(tt.opts)
			if len(got) != len(tt.expected) {
				t.Errorf("Filter() got %d packages, want %d", len(got), len(tt.expected))
			}
			for i, pkg := range got {
				if pkg != tt.expected[i] {
					t.Errorf("Filter()[%d] = %q, want %q", i, pkg, tt.expected[i])
				}
			}
		})
	}
}
