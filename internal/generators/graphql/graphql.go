// Package graphql generates the gqlgen-based GraphQL scaffold.
package graphql

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// Versions carries pinned GraphQL-related module versions.
type Versions struct {
	Gqlgen, FiberAdaptor string
}

// New returns the GraphQL generator pinned to v.
func New(v Versions) generator.Generator { return &gen{v: v} }

type gen struct{ v Versions }

func (gen) Name() string                 { return "graphql" }
func (gen) Applies(c config.Config) bool { return c.GraphQL }

func (g gen) Generate(ctx *generator.Context) error {
	files := []struct{ tmpl, out string }{
		{"graphql/schema.graphqls", "graph/schema.graphqls"},
		{"graphql/gqlgen.yml", "gqlgen.yml"},
		{"graphql/resolver.go", "graph/resolver.go"},
		{"graphql/handler.go", "internal/graphql/handler.go"},
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

	ctx.Manifest.Require("github.com/99designs/gqlgen", g.v.Gqlgen)
	if ctx.Config.Framework == config.FrameworkFiber {
		ctx.Manifest.Require("github.com/gofiber/adaptor/v2", g.v.FiberAdaptor)
	}
	return nil
}
