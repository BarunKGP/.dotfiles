package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Custom executes bash scripts/packages/<name>.sh with environment variables.
type Custom struct {
	dotfilesDir string
	osType      string
	packageMgr  string
}

func NewCustom(dotfilesDir, osType, packageMgr string) *Custom {
	return &Custom{
		dotfilesDir: dotfilesDir,
		osType:      osType,
		packageMgr:  packageMgr,
	}
}

// Execute runs the custom install script for a package.
func (c *Custom) Execute(pkg string, dryRun bool) error {
	scriptPath := filepath.Join(c.dotfilesDir, "scripts", "packages", pkg+".sh")

	// Check if script exists
	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf("custom script not found: %s", scriptPath)
	}

	args := []string{"bash", scriptPath}
	cmd := exec.Command(args[0], args[1:]...)

	// Set environment variables
	cmd.Env = append(os.Environ(),
		"OS_TYPE="+c.osType,
		"PACKAGE_MANAGER="+c.packageMgr,
		"DOTFILES="+c.dotfilesDir,
		"DRY_RUN="+(map[bool]string{true: "1", false: "0"}[dryRun]),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
