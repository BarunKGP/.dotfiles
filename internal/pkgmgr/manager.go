package pkgmgr

import (
	"os"
	"os/exec"

	"github.com/BarunKGP/dotfiles/dotctl/internal/osdetect"
)

// Manager is the interface for package managers.
type Manager interface {
	Install(pkg string, dryRun bool) error
	IsInstalled(pkg string) bool
	UpdateIndex(dryRun bool) error
	HasSudo() (bool, error)
	Name() string
}

type LookupFunc func(cmd string) (string, error)

var DefaultLookup LookupFunc = exec.LookPath

// Detect detects the available package manager.
func Detect(osType osdetect.OSType, lookup LookupFunc) Manager {
	if lookup == nil {
		lookup = DefaultLookup
	}

	managers := managerPriority(osType)
	for _, name := range managers {
		if pathExists(lookup, name) {
			return New(name, lookup)
		}
	}

	return nil
}

// New creates a new package manager by name.
func New(name string, lookup LookupFunc) Manager {
	if lookup == nil {
		lookup = DefaultLookup
	}

	switch name {
	case "apt":
		return NewApt(lookup)
	case "apk":
		return NewApk(lookup)
	case "brew":
		return NewBrew(lookup)
	case "dnf":
		return NewDnf(lookup)
	default:
		return nil
	}
}

// managerPriority returns the priority order of package managers for an OS.
func managerPriority(osType osdetect.OSType) []string {
	switch osType {
	case osdetect.MacOS:
		return []string{"brew"}
	case osdetect.WSL2:
		return []string{"apt", "brew"}
	case osdetect.Linux:
		return []string{"apt", "apk", "dnf", "pacman"}
	default:
		return []string{"apt", "apk", "brew", "dnf"}
	}
}

// pathExists checks both standard PATH and common system paths.
func pathExists(lookup LookupFunc, cmd string) bool {
	if _, err := lookup(cmd); err == nil {
		return true
	}
	// Also check common system paths for tools like apk, apt-get
	paths := []string{"/sbin/" + cmd, "/usr/sbin/" + cmd, "/usr/bin/" + cmd, "/bin/" + cmd}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// BaseManager provides common functionality.
type BaseManager struct {
	name   string
	lookup LookupFunc
}

func (m *BaseManager) HasSudo() (bool, error) {
	if os.Getuid() == 0 {
		return true, nil
	}
	_, err := m.lookup("sudo")
	return err == nil, nil
}

func (m *BaseManager) Name() string {
	return m.name
}

func (m *BaseManager) execCommand(args []string, dryRun bool) error {
	if dryRun {
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (m *BaseManager) commandExists(cmd string) bool {
	_, err := m.lookup(cmd)
	return err == nil
}
