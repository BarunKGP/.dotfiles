package manifest

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
		wantErr bool
	}{
		{
			name:    "empty",
			content: "",
			want:    0,
		},
		{
			name: "comments and blanks",
			content: `# Comment
zsh shell

# Another comment
bash shell`,
			want: 2,
		},
		{
			name: "with pm overrides",
			content: `neovim editor apk=neovim
espanso optional install=custom`,
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("Parse() returned nil")
				return
			}
			if len(got.Packages) != tt.want {
				t.Errorf("Parse() got %d packages, want %d", len(got.Packages), tt.want)
			}
		})
	}
}

func TestParseLineWithPMOverride(t *testing.T) {
	content := `neovim editor apk=neovim`
	manifest, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(manifest.Packages) != 1 {
		t.Fatalf("got %d packages, want 1", len(manifest.Packages))
	}

	pkg := manifest.Packages[0]
	if pkg.Canonical != "neovim" {
		t.Errorf("canonical = %q, want %q", pkg.Canonical, "neovim")
	}
	if pkg.Group != "editor" {
		t.Errorf("group = %q, want %q", pkg.Group, "editor")
	}
	if override, ok := pkg.PMOverrides["apk"]; !ok || override != "neovim" {
		t.Errorf("apk override = %q, %v; want %q, true", override, ok, "neovim")
	}
}

func TestParseLineWithCustomInstall(t *testing.T) {
	content := `espanso optional install=custom`
	manifest, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(manifest.Packages) != 1 {
		t.Fatalf("got %d packages, want 1", len(manifest.Packages))
	}

	pkg := manifest.Packages[0]
	if !pkg.CustomInstall {
		t.Errorf("CustomInstall = %v, want true", pkg.CustomInstall)
	}
}

func TestParseWhitespaceHandling(t *testing.T) {
	content := `  zsh   shell  `
	manifest, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(manifest.Packages) != 1 {
		t.Fatalf("got %d packages, want 1", len(manifest.Packages))
	}

	pkg := manifest.Packages[0]
	if pkg.Canonical != "zsh" {
		t.Errorf("canonical = %q, want %q", pkg.Canonical, "zsh")
	}
}
