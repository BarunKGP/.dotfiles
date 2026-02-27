package osdetect

import (
	"os"
	"strings"
)

type OSType string

const (
	WSL2  OSType = "wsl2"
	MacOS OSType = "macos"
	Linux OSType = "linux"
)

// Detector detects the OS type. It's injectable for testing.
type Detector struct {
	ProcVersionPath string
	UnameFunc       func() (string, error)
}

// NewDetector creates a new Detector with system defaults.
func NewDetector() *Detector {
	return &Detector{
		ProcVersionPath: "/proc/version",
		UnameFunc:       defaultUname,
	}
}

// Detect returns the detected OS type.
func (d *Detector) Detect() OSType {
	// Check for WSL2 first (most specific)
	if d.isWSL2() {
		return WSL2
	}

	// Check for macOS
	if d.isMacOS() {
		return MacOS
	}

	// Default to linux
	return Linux
}

func (d *Detector) isWSL2() bool {
	data, err := os.ReadFile(d.ProcVersionPath)
	if err != nil {
		return false
	}

	content := string(data)
	return strings.Contains(content, "microsoft") || strings.Contains(content, "WSL")
}

func (d *Detector) isMacOS() bool {
	// Check for macOS-specific paths
	_, err := os.Stat("/System/Library/CoreServices/Finder.app")
	return err == nil
}

// defaultUname returns the output of uname -s (not currently used but kept for future).
func defaultUname() (string, error) {
	// This would exec uname -s if needed
	return "", nil
}
