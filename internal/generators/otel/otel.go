// Package otel generates the OpenTelemetry scaffold (tracer provider + OTLP exporter).
package otel

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// Versions carries pinned OTel module versions.
type Versions struct {
	OTel, Contrib, Otelchi, Otelfiber string
}

// New returns the OTel generator pinned to v.
func New(v Versions) generator.Generator { return &gen{v: v} }

type gen struct{ v Versions }

func (gen) Name() string                 { return "otel" }
func (gen) Applies(c config.Config) bool { return c.OTel }

func (g gen) Generate(ctx *generator.Context) error {
	body, err := ctx.Renderer.Render("otel/telemetry.go", ctx.Config)
	if err != nil {
		return err
	}
	if err := ctx.Writer.Write("internal/telemetry/telemetry.go", body, 0o644); err != nil {
		return err
	}

	// Core SDK + OTLP/HTTP exporter.
	ctx.Manifest.Require("go.opentelemetry.io/otel", g.v.OTel)
	ctx.Manifest.Require("go.opentelemetry.io/otel/sdk", g.v.OTel)
	ctx.Manifest.Require("go.opentelemetry.io/otel/exporters/otlp/otlptrace", g.v.OTel)
	ctx.Manifest.Require("go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp", g.v.OTel)

	// Framework-specific HTTP instrumentation.
	switch ctx.Config.Framework {
	case config.FrameworkStdlib:
		ctx.Manifest.Require("go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp", g.v.Contrib)
	case config.FrameworkChi:
		ctx.Manifest.Require("github.com/riandyrn/otelchi", g.v.Otelchi)
	case config.FrameworkGin:
		ctx.Manifest.Require("go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin", g.v.Contrib)
	case config.FrameworkEcho:
		ctx.Manifest.Require("go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho", g.v.Contrib)
	case config.FrameworkFiber:
		ctx.Manifest.Require("github.com/gofiber/contrib/otelfiber", g.v.Otelfiber)
	}

	if ctx.Config.GRPC {
		ctx.Manifest.Require("go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc", g.v.Contrib)
	}
	return nil
}
