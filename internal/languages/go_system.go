package languages

import (
	"context"
	"os/exec"
	"strings"
)

type GoSystem struct{}

func NewGoSystem() Provider { return &GoSystem{} }

func (g *GoSystem) Name() string { return "go" }

func (g *GoSystem) Detect(ctx context.Context) bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func (g *GoSystem) Env(ctx context.Context) []EnvVar {
	goroot := g.goEnv("GOROOT")
	gopath := g.goEnv("GOPATH")
	if gopath == "" {
		gopath = "$HOME/go"
	}
	return []EnvVar{
		{Name: "GOROOT", Value: goroot, Quoted: true},
		{Name: "GOPATH", Value: gopath, Quoted: true},
		{Name: "PATH", Value: "$GOPATH/bin:$GOROOT/bin:$PATH", Quoted: false},
	}
}

func (g *GoSystem) Init(ctx context.Context) *Initialization { return nil }

func (g *GoSystem) goEnv(key string) string {
	out, err := exec.Command("go", "env", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
