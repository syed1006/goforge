// Package config defines the user-facing configuration consumed by every generator.
package config

import (
	"fmt"
	"regexp"
	"strings"
)

// Framework enumerates supported HTTP frameworks.
type Framework string

const (
	FrameworkStdlib Framework = "stdlib"
	FrameworkChi    Framework = "chi"
	FrameworkGin    Framework = "gin"
	FrameworkFiber  Framework = "fiber"
	FrameworkEcho   Framework = "echo"
)

// Frameworks lists every supported framework in display order.
var Frameworks = []Framework{
	FrameworkStdlib, FrameworkChi, FrameworkGin, FrameworkFiber, FrameworkEcho,
}

// Database enumerates supported database drivers.
type Database string

const (
	DatabaseNone     Database = "none"
	DatabasePostgres Database = "postgres"
	DatabaseMySQL    Database = "mysql"
	DatabaseSQLite   Database = "sqlite"
	DatabaseMongo    Database = "mongo"
	DatabaseRedis    Database = "redis"
)

// Databases lists every supported database driver in display order.
var Databases = []Database{
	DatabaseNone, DatabasePostgres, DatabaseMySQL, DatabaseSQLite, DatabaseMongo, DatabaseRedis,
}

// Config is the resolved scaffold configuration after prompts and flags are merged.
type Config struct {
	ProjectName string
	ModulePath  string
	OutputDir   string
	GoVersion   string

	Framework Framework
	Database  Database

	GRPC      bool
	GraphQL   bool
	HotReload bool
	Lint      bool
	Docker    bool
	CI        bool
}

var (
	projectNameRe = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	modulePathRe  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._/~-]*[a-zA-Z0-9]$`)
)

// Validate returns the first violation in c or nil if every field is acceptable.
func (c *Config) Validate() error {
	if !projectNameRe.MatchString(c.ProjectName) {
		return fmt.Errorf("project name %q must be lowercase, start with a letter, and contain only [a-z0-9_-]", c.ProjectName)
	}
	if c.ModulePath == "" || !modulePathRe.MatchString(c.ModulePath) || !strings.Contains(c.ModulePath, "/") {
		return fmt.Errorf("module path %q is not a valid Go module path", c.ModulePath)
	}
	if !containsFramework(Frameworks, c.Framework) {
		return fmt.Errorf("unknown framework %q", c.Framework)
	}
	if !containsDatabase(Databases, c.Database) {
		return fmt.Errorf("unknown database %q", c.Database)
	}
	if c.GoVersion == "" {
		return fmt.Errorf("go version is required")
	}
	return nil
}

// Slug returns a safe, lower-kebab-cased token derived from the project name.
func (c *Config) Slug() string {
	return strings.ToLower(strings.ReplaceAll(c.ProjectName, "_", "-"))
}

// BinaryName returns the binary name to use for generated Makefiles and Dockerfiles.
func (c *Config) BinaryName() string { return c.Slug() }

// ImportRoot returns the module path, suitable for prefixing internal imports.
func (c *Config) ImportRoot() string { return c.ModulePath }

func containsFramework(list []Framework, v Framework) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func containsDatabase(list []Database, v Database) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
