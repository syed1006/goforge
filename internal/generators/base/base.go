// Package base provides the framework-independent generator that writes the
// repository boilerplate every project shares.
package base

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the base generator.
func New() generator.Generator { return &gen{} }

type gen struct{}

func (gen) Name() string                 { return "base" }
func (gen) Applies(_ config.Config) bool { return true }

func (gen) Generate(ctx *generator.Context) error {
	data := struct {
		config.Config
		BinaryName string
	}{Config: ctx.Config, BinaryName: ctx.Config.BinaryName()}

	files := []struct {
		template string
		out      string
	}{
		{"base/README.md", "README.md"},
		{"base/Makefile", "Makefile"},
		{"base/gitignore", ".gitignore"},
		{"base/env", ".env.example"},
		{"base/buildinfo.go", "internal/buildinfo/buildinfo.go"},
		{"base/config.go", "internal/config/config.go"},
	}

	for _, f := range files {
		body, err := ctx.Renderer.Render(f.template, data)
		if err != nil {
			return err
		}
		if err := ctx.Writer.Write(f.out, body, 0o644); err != nil {
			return err
		}
	}
	return nil
}
