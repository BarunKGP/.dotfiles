package stow

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Stower manages stow operations.
type Stower struct {
	dotfilesDir string
	homeDir     string
}

// New creates a new Stower.
func New(dotfilesDir, homeDir string) *Stower {
	return &Stower{
		dotfilesDir: dotfilesDir,
		homeDir:     homeDir,
	}
}

// Check verifies that stow is available.
func Check() error {
	_, err := exec.LookPath("stow")
	if err != nil {
		return fmt.Errorf("stow not found in PATH")
	}
	return nil
}

// Apply stows a package.
func (s *Stower) Apply(pkg string, dryRun bool) error {
	args := []string{"stow"}

	if dryRun {
		args = append(args, "--no")
	}

	// Use --restow to replace existing links (idempotent)
	args = append(args, "--restow")

	// Use the package name as the source directory
	args = append(args, pkg)

	cmd := exec.Command("stow", args[1:]...)
	cmd.Dir = s.dotfilesDir
	cmd.Env = append(os.Environ(), fmt.Sprintf("HOME=%s", s.homeDir))

	// Capture stderr to provide better error messages
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		stderr_msg := stderr.String()
		if bytes.Contains(stderr.Bytes(), []byte("Conflict")) {
			return fmt.Errorf("stow conflict for %s: file already exists. Remove it, then retry", pkg)
		}
		if stderr_msg != "" {
			return fmt.Errorf("stow failed for %s: %s", pkg, stderr_msg)
		}
		return fmt.Errorf("stow failed for %s: %w", pkg, err)
	}

	return nil
}

// GetDefaultPackages returns the default packages to stow.
func GetDefaultPackages(minimal bool) []string {
	packages := []string{"bash", "git", "nvim"}

	if !minimal {
		packages = append(packages, "zsh", "tmux")
	}

	return packages
}
