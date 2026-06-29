// Package grpc generates the gRPC scaffold: server wrapper, buf config, sample proto.
package grpc

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// New returns the gRPC generator pinned to grpcVersion.
func New(grpcVersion string) generator.Generator { return &gen{version: grpcVersion} }

type gen struct{ version string }

func (gen) Name() string                 { return "grpc" }
func (gen) Applies(c config.Config) bool { return c.GRPC }

func (g gen) Generate(ctx *generator.Context) error {
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

	ctx.Manifest.Require("google.golang.org/grpc", g.version)
	return nil
}
