package osdetect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectWSL2WithMicrosoftString(t *testing.T) {
	tmpDir := t.TempDir()
	procVersionPath := filepath.Join(tmpDir, "proc_version")
	os.WriteFile(procVersionPath, []byte("4.19.128-microsoft-standard"), 0644)

	detector := &Detector{
		ProcVersionPath: procVersionPath,
	}

	if got := detector.Detect(); got != WSL2 {
		t.Errorf("Detect() = %v, want %v", got, WSL2)
	}
}

func TestDetectWSL2WithWSLString(t *testing.T) {
	tmpDir := t.TempDir()
	procVersionPath := filepath.Join(tmpDir, "proc_version")
	os.WriteFile(procVersionPath, []byte("4.19.128 WSL (2)"), 0644)

	detector := &Detector{
		ProcVersionPath: procVersionPath,
	}

	if got := detector.Detect(); got != WSL2 {
		t.Errorf("Detect() = %v, want %v", got, WSL2)
	}
}

func TestDetectLinuxPlain(t *testing.T) {
	tmpDir := t.TempDir()
	procVersionPath := filepath.Join(tmpDir, "proc_version")
	os.WriteFile(procVersionPath, []byte("5.10.0-generic"), 0644)

	detector := &Detector{
		ProcVersionPath: procVersionPath,
	}

	if got := detector.Detect(); got != Linux {
		t.Errorf("Detect() = %v, want %v", got, Linux)
	}
}

func TestDetectNoProcVersion(t *testing.T) {
	detector := &Detector{
		ProcVersionPath: "/nonexistent/proc/version",
	}

	if got := detector.Detect(); got != Linux {
		t.Errorf("Detect() = %v, want %v", got, Linux)
	}
}
