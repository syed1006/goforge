// Package generator defines the contract every scaffold feature implements
// and the shared context the orchestrator hands them at runtime.
package generator

import (
	"io/fs"

	"github.com/syed1006/goforge/internal/config"
)

// Generator produces a slice of the final project tree based on the resolved config.
//
// Implementations should be stateless and reentrant — the orchestrator may run
// generators in any order it likes.
type Generator interface {
	// Name returns a short, human-readable identifier (used in logs).
	Name() string
	// Applies reports whether the generator should run for the given config.
	Applies(cfg config.Config) bool
	// Generate produces files via ctx.Writer and registers module deps via ctx.Manifest.
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

// Context is handed to every generator and bundles the resolved config with
// the rendering and writing capabilities the generator needs.
type Context struct {
	Config   config.Config
	Renderer Renderer
	Writer   Writer
	Manifest *Manifest
}
