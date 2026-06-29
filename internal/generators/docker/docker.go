// Package docker generates Dockerfile, .dockerignore, and docker-compose.yml.
package docker

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the Docker generator.
func New() generator.Generator { return &gen{} }

type gen struct{}

func (gen) Name() string                 { return "docker" }
func (gen) Applies(c config.Config) bool { return c.Docker }

func (gen) Generate(ctx *generator.Context) error {
	data := struct {
		config.Config
		BinaryName string
	}{Config: ctx.Config, BinaryName: ctx.Config.BinaryName()}

	files := []struct{ tmpl, out string }{
		{"docker/Dockerfile", "Dockerfile"},
		{"docker/dockerignore", ".dockerignore"},
		{"docker/compose.yml", "docker-compose.yml"},
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
