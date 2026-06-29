// Package framework holds the HTTP framework generators.
//
// Each supported framework owns the entrypoint (cmd/<bin>/main.go) and the
// http server wiring (internal/server). Exactly one of these generators runs
// per scaffold, selected by the cfg.Framework value.
package framework

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// NewStdlib returns the net/http-only framework generator.
func NewStdlib() generator.Generator { return &stdlib{} }

type stdlib struct{}

func (stdlib) Name() string                       { return "framework/stdlib" }
func (stdlib) Applies(c config.Config) bool       { return c.Framework == config.FrameworkStdlib }

func (s stdlib) Generate(ctx *generator.Context) error {
	data := struct {
		config.Config
		BinaryName string
	}{Config: ctx.Config, BinaryName: ctx.Config.BinaryName()}

	files := []struct{ tmpl, out string }{
		{"framework/stdlib/main.go", "cmd/" + ctx.Config.BinaryName() + "/main.go"},
		{"framework/stdlib/server.go", "internal/server/server.go"},
	}
	for _, f := range files {
		body, err := ctx.Renderer.Render(f.tmpl, data)
		if err != nil {
			return err
		}
		if err := ctx.Writer.Write(f.out, body, 0o644); err != nil {
			return err
		}
	}
	return nil
}
