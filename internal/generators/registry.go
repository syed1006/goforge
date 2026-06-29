// Package generators wires every concrete generator into the default Registry.
package generators

import (
	"github.com/syed1006/goforge/internal/generator"
	"github.com/syed1006/goforge/internal/generators/base"
	"github.com/syed1006/goforge/internal/generators/database"
	"github.com/syed1006/goforge/internal/generators/framework"
	"github.com/syed1006/goforge/internal/generators/graphql"
	grpcgen "github.com/syed1006/goforge/internal/generators/grpc"
)

// Default returns the default generator registry, with generators registered
// in the order they should run.
func Default() *generator.Registry {
	reg := generator.NewRegistry()
	reg.Register(base.New())
	reg.Register(framework.All()...)
	reg.Register(database.All()...)
	reg.Register(grpcgen.New())
	reg.Register(graphql.New())
	return reg
}
