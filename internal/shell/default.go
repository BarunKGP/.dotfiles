package shell

import (
	"fmt"
	"os"
	"os/exec"
)

// SetDefaultResult is the result of setting the default shell.
type SetDefaultResult struct {
	Success bool
	Shell   string
	Error   error
}

// SetDefault attempts to set the default shell via chsh.
// Returns non-fatal error if chsh is unavailable or user not in /etc/passwd.
func SetDefault(name string) *SetDefaultResult {
	// Find full path to shell
	shellPath, err := exec.LookPath(name)
	if err != nil {
		return &SetDefaultResult{
			Success: false,
			Shell:   name,
			Error:   fmt.Errorf("shell not found in PATH: %s", name),
		}
	}

	// Try chsh
	if _, err := exec.LookPath("chsh"); err != nil {
		return &SetDefaultResult{
			Success: false,
			Shell:   name,
			Error:   fmt.Errorf("chsh not available"),
		}
	}

	cmd := exec.Command("chsh", "-s", shellPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return &SetDefaultResult{
			Success: false,
			Shell:   name,
			Error:   fmt.Errorf("chsh failed: %w", err),
		}
	}

	return &SetDefaultResult{
		Success: true,
		Shell:   name,
		Error:   nil,
	}
}
