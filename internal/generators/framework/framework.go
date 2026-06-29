// Package framework generates HTTP framework scaffolds; exactly one runs per project.
package framework

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// All returns every framework generator (one per supported framework).
func All() []generator.Generator {
	return []generator.Generator{
		newImpl(config.FrameworkStdlib, "stdlib", nil),
		newImpl(config.FrameworkChi, "chi", []moduleDep{
			{"github.com/go-chi/chi/v5", "latest"},
		}),
		newImpl(config.FrameworkGin, "gin", []moduleDep{
			{"github.com/gin-gonic/gin", "latest"},
		}),
		newImpl(config.FrameworkFiber, "fiber", []moduleDep{
			{"github.com/gofiber/fiber/v2", "latest"},
		}),
		newImpl(config.FrameworkEcho, "echo", []moduleDep{
			{"github.com/labstack/echo/v4", "latest"},
		}),
	}
}

type moduleDep struct {
	module  string
	version string
}

type impl struct {
	framework config.Framework
	dir       string
	deps      []moduleDep
}

func newImpl(fw config.Framework, dir string, deps []moduleDep) *impl {
	return &impl{framework: fw, dir: dir, deps: deps}
}

func (i *impl) Name() string                 { return "framework/" + i.dir }
func (i *impl) Applies(c config.Config) bool { return c.Framework == i.framework }

func (i *impl) Generate(ctx *generator.Context) error {
	data := struct {
		config.Config
		BinaryName string
	}{Config: ctx.Config, BinaryName: ctx.Config.BinaryName()}

	files := []struct{ tmpl, out string }{
		{"framework/" + i.dir + "/main.go", "cmd/" + ctx.Config.BinaryName() + "/main.go"},
		{"framework/" + i.dir + "/server.go", "internal/server/server.go"},
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

	for _, d := range i.deps {
		ctx.Manifest.Require(d.module, d.version)
	}
	return nil
}
