// Package graphql holds the gqlgen-based GraphQL generator.
//
// Produces:
//   - graph/schema.graphqls            — sample schema
//   - gqlgen.yml                       — gqlgen config
//   - graph/resolver.go                — Resolver root (handwritten; gqlgen
//                                        will not overwrite it)
//   - internal/graphql/handler.go      — framework-agnostic HTTP handler
//
// The orchestrator runs `gqlgen generate` as a post-step (see internal/scaffold)
// when c.GraphQL is true so the generated `graph/generated.go` and
// `graph/schema.resolvers.go` exist before the project is first built.
package graphql

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the GraphQL generator.
func New() generator.Generator { return &gen{} }

type gen struct{}

func (gen) Name() string                 { return "graphql" }
func (gen) Applies(c config.Config) bool { return c.GraphQL }

func (gen) Generate(ctx *generator.Context) error {
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

	ctx.Manifest.Require("github.com/99designs/gqlgen", "latest")
	if ctx.Config.Framework == config.FrameworkFiber {
		ctx.Manifest.Require("github.com/gofiber/adaptor/v2", "latest")
	}
	return nil
}
