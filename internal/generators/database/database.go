// Package database holds the storage driver generators.
//
// Exactly one of these generators runs per scaffold (selected by cfg.Database),
// and each writes internal/storage/storage.go plus the module deps it needs.
package database

import (
	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

// All returns every database generator (one per supported driver, plus a noop
// for cfg.Database == none — implemented by simply having Applies() return false).
func All() []generator.Generator {
	return []generator.Generator{
		newImpl(config.DatabasePostgres, "postgres", []moduleDep{{"github.com/jackc/pgx/v5", "latest"}}),
		newImpl(config.DatabaseMySQL, "mysql", []moduleDep{{"github.com/go-sql-driver/mysql", "latest"}}),
		newImpl(config.DatabaseSQLite, "sqlite", []moduleDep{{"modernc.org/sqlite", "latest"}}),
		newImpl(config.DatabaseMongo, "mongo", []moduleDep{{"go.mongodb.org/mongo-driver", "latest"}}),
		newImpl(config.DatabaseRedis, "redis", []moduleDep{{"github.com/redis/go-redis/v9", "latest"}}),
	}
}

type moduleDep struct {
	module  string
	version string
}

type impl struct {
	driver config.Database
	dir    string
	deps   []moduleDep
}

func newImpl(d config.Database, dir string, deps []moduleDep) *impl {
	return &impl{driver: d, dir: dir, deps: deps}
}

func (i *impl) Name() string                 { return "database/" + i.dir }
func (i *impl) Applies(c config.Config) bool { return c.Database == i.driver }

func (i *impl) Generate(ctx *generator.Context) error {
	body, err := ctx.Renderer.Render("database/"+i.dir+"/storage.go", ctx.Config)
	if err != nil {
		return err
	}
	if err := ctx.Writer.Write("internal/storage/storage.go", body, 0o644); err != nil {
		return err
	}
	for _, d := range i.deps {
		ctx.Manifest.Require(d.module, d.version)
	}
	return nil
}
