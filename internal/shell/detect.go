package shell

import (
	"os/exec"
)

// LookupFunc is a function that finds a command in PATH.
type LookupFunc func(cmd string) (string, error)

// DefaultLookup uses exec.LookPath.
var DefaultLookup LookupFunc = exec.LookPath

// AvailableShells returns the list of available shells.
func AvailableShells(lookup LookupFunc) []string {
	if lookup == nil {
		lookup = DefaultLookup
	}

	shells := []string{"bash", "zsh", "fish", "sh"}
	var available []string

	for _, shell := range shells {
		if _, err := lookup(shell); err == nil {
			available = append(available, shell)
		}
	}

	return available
}
