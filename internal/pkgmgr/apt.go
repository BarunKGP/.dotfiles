package pkgmgr

import (
	"os"
	"os/exec"
	"sync"
)

type Apt struct {
	BaseManager
	indexUpdated sync.Once
	indexErr     error
}

func NewApt(lookup LookupFunc) *Apt {
	return &Apt{
		BaseManager: BaseManager{
			name:   "apt",
			lookup: lookup,
		},
	}
}

func (a *Apt) Install(pkg string, dryRun bool) error {
	if dryRun {
		return nil
	}

	// Ensure index is updated once
	if err := a.UpdateIndex(dryRun); err != nil {
		return err
	}

	sudoNeeded, err := a.HasSudo()
	if err != nil {
		return err
	}

	args := []string{"apt-get", "install", "-y", pkg}
	if sudoNeeded && os.Getuid() != 0 {
		args = append([]string{"sudo"}, args...)
	}

	cmd := a.buildCommand(args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (a *Apt) IsInstalled(pkg string) bool {
	return a.commandExists(pkg)
}

func (a *Apt) UpdateIndex(dryRun bool) error {
	var err error
	a.indexUpdated.Do(func() {
		if dryRun {
			return
		}

		sudoNeeded, sudoErr := a.HasSudo()
		if sudoErr != nil {
			err = sudoErr
			return
		}

		args := []string{"apt-get", "update"}
		if sudoNeeded && os.Getuid() != 0 {
			args = append([]string{"sudo"}, args...)
		}

		cmd := a.buildCommand(args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	})
	return err
}

func (a *Apt) buildCommand(args []string) *exec.Cmd {
	return exec.Command(args[0], args[1:]...)
}
