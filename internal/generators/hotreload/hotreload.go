// Package hotreload generates the .air.toml configuration.
package hotreload

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the hot-reload generator.
func New() generator.Generator { return &gen{} }

type gen struct{}

func (gen) Name() string                 { return "hotreload" }
func (gen) Applies(c config.Config) bool { return c.HotReload }

func (gen) Generate(ctx *generator.Context) error {
	data := struct {
		config.Config
		BinaryName string
	}{Config: ctx.Config, BinaryName: ctx.Config.BinaryName()}

	body, err := ctx.Renderer.Render("hotreload/air.toml", data)
	if err != nil {
		return err
	}
	if err := ctx.Writer.Write(".air.toml", body, 0o644); err != nil {
		return err
	}
	ctx.Manifest.Tool("air", "go install github.com/air-verse/air@latest")
	return nil
}
