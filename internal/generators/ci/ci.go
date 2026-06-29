// Package ci writes the GitHub Actions workflow for the generated project.
package ci

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the CI generator.
func New() generator.Generator { return &gen{} }

type gen struct{}

func (gen) Name() string                 { return "ci" }
func (gen) Applies(c config.Config) bool { return c.CI }

func (gen) Generate(ctx *generator.Context) error {
	data := struct {
		config.Config
		BinaryName string
	}{Config: ctx.Config, BinaryName: ctx.Config.BinaryName()}
	body, err := ctx.Renderer.Render("ci/ci.yml", data)
	if err != nil {
		return err
	}
	return ctx.Writer.Write(".github/workflows/ci.yml", body, 0o644)
}
