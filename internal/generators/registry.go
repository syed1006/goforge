// Package generators wires every concrete generator into the default Registry.
package generators

import (
	"github.com/syed1006/goforge/internal/generator"
	"github.com/syed1006/goforge/internal/generators/base"
	"github.com/syed1006/goforge/internal/generators/framework"
)

// Default returns the default generator registry, with generators registered
// in the order they should run.
func Default() *generator.Registry {
	reg := generator.NewRegistry()
	reg.Register(
		base.New(),
		framework.NewStdlib(),
	)
	return reg
}
