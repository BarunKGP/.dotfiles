package pkgmgr

import (
	"os"
	"os/exec"
)

type Dnf struct {
	BaseManager
}

func NewDnf(lookup LookupFunc) *Dnf {
	return &Dnf{
		BaseManager: BaseManager{
			name:   "dnf",
			lookup: lookup,
		},
	}
}

func (d *Dnf) Install(pkg string, dryRun bool) error {
	if dryRun {
		return nil
	}

	sudoNeeded, err := d.HasSudo()
	if err != nil {
		return err
	}

	args := []string{"dnf", "install", "-y", pkg}
	if sudoNeeded && os.Getuid() != 0 {
		args = append([]string{"sudo"}, args...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (d *Dnf) IsInstalled(pkg string) bool {
	return d.commandExists(pkg)
}

func (d *Dnf) UpdateIndex(dryRun bool) error {
	// Dnf doesn't require explicit index update
	return nil
}
