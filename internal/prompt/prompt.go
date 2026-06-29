// Package prompt drives the interactive flow that gathers a scaffold Config from
// the user. Only fields left empty in the seed are asked about; everything else
// (typically supplied via flags) is passed through.
package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/syed1006/goforge/internal/config"
)

// Ask presents an interactive form, seeded by the supplied config, and returns
// the merged result. Fields already set in the seed are not re-asked.
func Ask(seed config.Config) (config.Config, error) {
	cfg := seed
	if cfg.GoVersion == "" {
		cfg.GoVersion = defaultGoVersion()
	}

	groups := []*huh.Group{}

	if cfg.ProjectName == "" || cfg.ModulePath == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("Lowercase, [a-z0-9_-], starts with a letter").
				Value(&cfg.ProjectName).
				Validate(validateProjectName),
			huh.NewInput().
				Title("Module path").
				Description("e.g. github.com/you/myapi").
				Value(&cfg.ModulePath).
				Validate(validateModulePath),
		))
	}

	if cfg.Framework == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[config.Framework]().
				Title("HTTP framework").
				Options(frameworkOptions()...).
				Value(&cfg.Framework),
		))
	}

	if cfg.Database == "" {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[config.Database]().
				Title("Database driver").
				Options(databaseOptions()...).
				Value(&cfg.Database),
		))
	}

	groups = append(groups, huh.NewGroup(
		huh.NewConfirm().Title("gRPC server?").Value(&cfg.GRPC),
		huh.NewConfirm().Title("GraphQL server (gqlgen)?").Value(&cfg.GraphQL),
		huh.NewConfirm().Title("OpenTelemetry tracing (OTLP)?").Value(&cfg.OTel),
		huh.NewConfirm().Title("Prometheus /metrics endpoint?").Value(&cfg.Metrics),
		huh.NewConfirm().Title("Hot reload (air)?").Value(&cfg.HotReload),
		huh.NewConfirm().Title("Linting (golangci-lint + pre-commit)?").Value(&cfg.Lint),
		huh.NewConfirm().Title("Docker + docker-compose?").Value(&cfg.Docker),
		huh.NewConfirm().Title("GitHub Actions CI?").Value(&cfg.CI),
	))

	if err := huh.NewForm(groups...).Run(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func frameworkOptions() []huh.Option[config.Framework] {
	out := make([]huh.Option[config.Framework], 0, len(config.Frameworks))
	for _, f := range config.Frameworks {
		out = append(out, huh.NewOption(string(f), f))
	}
	return out
}

func databaseOptions() []huh.Option[config.Database] {
	out := make([]huh.Option[config.Database], 0, len(config.Databases))
	for _, d := range config.Databases {
		out = append(out, huh.NewOption(string(d), d))
	}
	return out
}

func validateProjectName(s string) error {
	c := config.Config{ProjectName: s, ModulePath: "github.com/x/y", GoVersion: "1.23", Framework: config.FrameworkStdlib, Database: config.DatabaseNone}
	if err := c.Validate(); err != nil && strings.Contains(err.Error(), "project name") {
		return err
	}
	return nil
}

func validateModulePath(s string) error {
	c := config.Config{ProjectName: "x", ModulePath: s, GoVersion: "1.23", Framework: config.FrameworkStdlib, Database: config.DatabaseNone}
	if err := c.Validate(); err != nil && strings.Contains(err.Error(), "module path") {
		return err
	}
	return nil
}

func defaultGoVersion() string { return "1.23" }

// Summary returns a short, plain-text summary of the resolved config — useful
// for confirmation prompts and dry-run output.
func Summary(cfg config.Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Project:    %s\n", cfg.ProjectName)
	fmt.Fprintf(&b, "Module:     %s\n", cfg.ModulePath)
	fmt.Fprintf(&b, "Go:         %s\n", cfg.GoVersion)
	fmt.Fprintf(&b, "Framework:  %s\n", cfg.Framework)
	fmt.Fprintf(&b, "Database:   %s\n", cfg.Database)
	fmt.Fprintf(&b, "gRPC:       %v\n", cfg.GRPC)
	fmt.Fprintf(&b, "GraphQL:    %v\n", cfg.GraphQL)
	fmt.Fprintf(&b, "OTel:       %v\n", cfg.OTel)
	fmt.Fprintf(&b, "Metrics:    %v\n", cfg.Metrics)
	fmt.Fprintf(&b, "Hot reload: %v\n", cfg.HotReload)
	fmt.Fprintf(&b, "Lint:       %v\n", cfg.Lint)
	fmt.Fprintf(&b, "Docker:     %v\n", cfg.Docker)
	fmt.Fprintf(&b, "CI:         %v\n", cfg.CI)
	return b.String()
}
