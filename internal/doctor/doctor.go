package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Report struct {
	Name    string
	Status  string // "PASS", "FAIL", "WARN"
	Details string
}

// Diagnose runs diagnostic checks on the dotfiles setup.
func Diagnose(dotfilesDir string) []Report {
	var reports []Report

	// Check 1: stow in PATH
	reports = append(reports, checkStowInPath())

	// Check 2: git in PATH
	reports = append(reports, checkGitInPath())

	// Check 3: packages.conf exists
	reports = append(reports, checkPackagesConf(dotfilesDir))

	// Check 4: ~/.config/shell/local.sh exists
	reports = append(reports, checkLocalShellConfig())

	// Check 5: ~/.bashrc is a symlink
	reports = append(reports, checkBashrcSymlink())

	return reports
}

func checkStowInPath() Report {
	_, err := exec.LookPath("stow")
	if err == nil {
		return Report{
			Name:    "stow in PATH",
			Status:  "PASS",
			Details: "stow command found",
		}
	}
	return Report{
		Name:    "stow in PATH",
		Status:  "FAIL",
		Details: "stow command not found in PATH",
	}
}

func checkGitInPath() Report {
	_, err := exec.LookPath("git")
	if err == nil {
		return Report{
			Name:    "git in PATH",
			Status:  "PASS",
			Details: "git command found",
		}
	}
	return Report{
		Name:    "git in PATH",
		Status:  "FAIL",
		Details: "git command not found in PATH",
	}
}

func checkPackagesConf(dotfilesDir string) Report {
	configPath := filepath.Join(dotfilesDir, "packages.conf")
	_, err := os.Stat(configPath)
	if err == nil {
		return Report{
			Name:    "packages.conf exists",
			Status:  "PASS",
			Details: fmt.Sprintf("Found at %s", configPath),
		}
	}
	return Report{
		Name:    "packages.conf exists",
		Status:  "FAIL",
		Details: fmt.Sprintf("Not found at %s", configPath),
	}
}

func checkLocalShellConfig() Report {
	home, err := os.UserHomeDir()
	if err != nil {
		return Report{
			Name:    "~/.config/shell/local.sh exists",
			Status:  "WARN",
			Details: "Could not determine home directory",
		}
	}

	configPath := filepath.Join(home, ".config", "shell", "local.sh")
	_, err = os.Stat(configPath)
	if err == nil {
		return Report{
			Name:    "~/.config/shell/local.sh exists",
			Status:  "PASS",
			Details: fmt.Sprintf("Found at %s", configPath),
		}
	}
	return Report{
		Name:    "~/.config/shell/local.sh exists",
		Status:  "WARN",
		Details: "Not yet created (this is OK)",
	}
}

func checkBashrcSymlink() Report {
	home, err := os.UserHomeDir()
	if err != nil {
		return Report{
			Name:    "~/.bashrc is a symlink",
			Status:  "WARN",
			Details: "Could not determine home directory",
		}
	}

	bashrcPath := filepath.Join(home, ".bashrc")
	fi, err := os.Lstat(bashrcPath)
	if err != nil {
		return Report{
			Name:    "~/.bashrc is a symlink",
			Status:  "WARN",
			Details: "File not found (this is OK if you use zsh)",
		}
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return Report{
			Name:    "~/.bashrc is a symlink",
			Status:  "PASS",
			Details: fmt.Sprintf("Symlink found at %s", bashrcPath),
		}
	}
	return Report{
		Name:    "~/.bashrc is a symlink",
		Status:  "WARN",
		Details: ".bashrc exists but is not a symlink",
	}
}
