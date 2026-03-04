package languages

import (
	"context"
	"os"
	"path/filepath"
)

type NodeNvm struct{}

func NewNodeNvm() Provider { return &NodeNvm{} }

func (n *NodeNvm) Name() string { return "node" }

func (n *NodeNvm) Detect(ctx context.Context) bool {
	home, _ := os.UserHomeDir()
	_, err := os.Stat(filepath.Join(home, ".nvm", "nvm.sh"))
	return err == nil
}

func (n *NodeNvm) Env(ctx context.Context) []EnvVar {
	return []EnvVar{
		{Name: "NVM_DIR", Value: "$HOME/.nvm", Quoted: true},
	}
}

func (n *NodeNvm) Init(ctx context.Context) *Initialization {
	return &Initialization{Snippet: `[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"`}
}
