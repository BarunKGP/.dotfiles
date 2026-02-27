package pkgmgr

import (
	"os"
	"os/exec"
)

type Apk struct {
	BaseManager
}

func NewApk(lookup LookupFunc) *Apk {
	return &Apk{
		BaseManager: BaseManager{
			name:   "apk",
			lookup: lookup,
		},
	}
}

func (a *Apk) Install(pkg string, dryRun bool) error {
	if dryRun {
		return nil
	}

	sudoNeeded, err := a.HasSudo()
	if err != nil {
		return err
	}

	args := []string{"apk", "add", pkg}
	if sudoNeeded && os.Getuid() != 0 {
		args = append([]string{"sudo"}, args...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (a *Apk) IsInstalled(pkg string) bool {
	return a.commandExists(pkg)
}

func (a *Apk) UpdateIndex(dryRun bool) error {
	// Alpine doesn't require explicit index update
	return nil
}
