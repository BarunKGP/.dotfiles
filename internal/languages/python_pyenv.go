package languages

import (
	"context"
	"os"
	"path/filepath"
)

type PythonPyenv struct{}

func NewPythonPyenv() Provider { return &PythonPyenv{} }

func (p *PythonPyenv) Name() string { return "python" }

func (p *PythonPyenv) Detect(ctx context.Context) bool {
	home, _ := os.UserHomeDir()
	_, err := os.Stat(filepath.Join(home, ".pyenv", "bin", "pyenv"))
	return err == nil
}

func (p *PythonPyenv) Env(ctx context.Context) []EnvVar {
	return []EnvVar{
		{Name: "PYENV_ROOT", Value: "$HOME/.pyenv", Quoted: true},
		{Name: "PATH", Value: "$PYENV_ROOT/bin:$PATH", Quoted: false},
	}
}

func (p *PythonPyenv) Init(ctx context.Context) *Initialization {
	return &Initialization{Snippet: `eval "$(pyenv init --path)" 2>/dev/null || true`}
}
