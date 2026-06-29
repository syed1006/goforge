// Package metrics generates the Prometheus /metrics scaffold.
package metrics

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// Versions carries pinned metrics module versions.
type Versions struct {
	Prom, FiberAdaptor string
}

// New returns the metrics generator pinned to v.
func New(v Versions) generator.Generator { return &gen{v: v} }

type gen struct{ v Versions }

func (gen) Name() string                 { return "metrics" }
func (gen) Applies(c config.Config) bool { return c.Metrics }

func (g gen) Generate(ctx *generator.Context) error {
	body, err := ctx.Renderer.Render("metrics/metrics.go", ctx.Config)
	if err != nil {
		return err
	}
	if err := ctx.Writer.Write("internal/metrics/metrics.go", body, 0o644); err != nil {
		return err
	}
	ctx.Manifest.Require("github.com/prometheus/client_golang", g.v.Prom)
	if ctx.Config.Framework == config.FrameworkFiber {
		ctx.Manifest.Require("github.com/gofiber/adaptor/v2", g.v.FiberAdaptor)
	}
	return nil
}
