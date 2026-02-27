package main

import (
	"os"

	"github.com/BarunKGP/dotfiles/dotctl/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
