// Package generators wires every concrete generator into the default Registry.
package generators

import (
	"github.com/syed1006/goforge/internal/generator"
	"github.com/syed1006/goforge/internal/generators/base"
	"github.com/syed1006/goforge/internal/generators/ci"
	"github.com/syed1006/goforge/internal/generators/database"
	"github.com/syed1006/goforge/internal/generators/docker"
	"github.com/syed1006/goforge/internal/generators/framework"
	"github.com/syed1006/goforge/internal/generators/graphql"
	grpcgen "github.com/syed1006/goforge/internal/generators/grpc"
	"github.com/syed1006/goforge/internal/generators/hotreload"
	"github.com/syed1006/goforge/internal/generators/lint"
)

// Default returns the default generator registry using pinned module versions.
func Default() *generator.Registry {
	reg := generator.NewRegistry()
	reg.Register(base.New())
	reg.Register(framework.All(framework.Versions{
		Chi: VersionChi, Gin: VersionGin, Fiber: VersionFiber, Echo: VersionEcho,
	})...)
	reg.Register(database.All(database.Versions{
		Postgres: VersionPgx, MySQL: VersionMySQL, SQLite: VersionSQLite,
		Mongo: VersionMongo, Redis: VersionRedis,
	})...)
	reg.Register(grpcgen.New(VersionGRPC))
	reg.Register(graphql.New(graphql.Versions{Gqlgen: VersionGqlgen, FiberAdaptor: VersionFiberAdaptor}))
	reg.Register(hotreload.New())
	reg.Register(lint.New())
	reg.Register(docker.New())
	reg.Register(ci.New())
	return reg
}
