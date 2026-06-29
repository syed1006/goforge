// Package grpc owns the gRPC scaffold.
//
// Produces:
//   - internal/grpcsrv/server.go   — wrapper around *grpc.Server with health+reflection
//   - proto/<service>/v1/<service>.proto — sample protobuf file
//   - buf.yaml, buf.gen.yaml      — buf configuration
//
// Runtime deps registered on the manifest: google.golang.org/grpc.
package grpc

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the gRPC generator.
func New() generator.Generator { return &gen{} }

type gen struct{}

func (gen) Name() string                 { return "grpc" }
func (gen) Applies(c config.Config) bool { return c.GRPC }

func (gen) Generate(ctx *generator.Context) error {
	cfg := ctx.Config
	protoDir := "proto/" + ctx.Config.Slug() + "/v1"

	files := []struct{ tmpl, out string }{
		{"grpc/server.go", "internal/grpcsrv/server.go"},
		{"grpc/buf.yaml", "buf.yaml"},
		{"grpc/buf.gen.yaml", "buf.gen.yaml"},
		{"grpc/ping.proto", protoDir + "/ping.proto"},
	}
	for _, f := range files {
		body, err := ctx.Renderer.Render(f.tmpl, cfg)
		if err != nil {
			return err
		}
		if err := ctx.Writer.Write(f.out, body, 0o644); err != nil {
			return err
		}
	}

	ctx.Manifest.Require("google.golang.org/grpc", "latest")
	return nil
}
