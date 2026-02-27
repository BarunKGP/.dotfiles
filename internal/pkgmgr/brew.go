package pkgmgr

import (
	"os"
	"os/exec"
)

type Brew struct {
	BaseManager
}

func NewBrew(lookup LookupFunc) *Brew {
	return &Brew{
		BaseManager: BaseManager{
			name:   "brew",
			lookup: lookup,
		},
	}
}

func (b *Brew) Install(pkg string, dryRun bool) error {
	if dryRun {
		return nil
	}

	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (b *Brew) IsInstalled(pkg string) bool {
	return b.commandExists(pkg)
}

func (b *Brew) UpdateIndex(dryRun bool) error {
	// Brew updates indices automatically
	return nil
}
