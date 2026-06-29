// Package lint writes the linting and formatting configuration.
//
// Produces:
//   - .golangci.yml         — golangci-lint config with sensible defaults
//   - .pre-commit-config.yaml — pre-commit hooks (fmt, mod-tidy, golangci-lint)
//   - .editorconfig         — common editor settings
package lint

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the lint generator.
func New() generator.Generator { return &gen{} }

type gen struct{}

func (gen) Name() string                 { return "lint" }
func (gen) Applies(c config.Config) bool { return c.Lint }

func (gen) Generate(ctx *generator.Context) error {
	files := []struct{ tmpl, out string }{
		{"lint/golangci.yml", ".golangci.yml"},
		{"lint/pre-commit.yml", ".pre-commit-config.yaml"},
		{"lint/editorconfig", ".editorconfig"},
	}
	for _, f := range files {
		body, err := ctx.Renderer.Render(f.tmpl, ctx.Config)
		if err != nil {
			return err
		}
		if err := ctx.Writer.Write(f.out, body, 0o644); err != nil {
			return err
		}
	}
	ctx.Manifest.Tool("golangci-lint", "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	ctx.Manifest.Tool("pre-commit", "pipx install pre-commit  # then: pre-commit install")
	return nil
}
