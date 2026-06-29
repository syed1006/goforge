// Package generator defines the contract every scaffold feature implements.
package generator

import (
	"io/fs"

	"github.com/syed1006/goforge/internal/config"
)

// Generator produces part of the project tree.
type Generator interface {
	Name() string
	Applies(cfg config.Config) bool
	Generate(ctx *Context) error
}

// Renderer renders a named template against the supplied data.
type Renderer interface {
	Render(name string, data any) ([]byte, error)
}

// Writer commits a file relative to the project root.
type Writer interface {
	Write(relPath string, content []byte, mode fs.FileMode) error
}

// Context bundles everything a Generator needs from the orchestrator.
type Context struct {
	Config   config.Config
	Renderer Renderer
	Writer   Writer
	Manifest *Manifest
}
