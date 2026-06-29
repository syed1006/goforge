package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generators"
	"github.com/syed1006/goforge/internal/prompt"
	"github.com/syed1006/goforge/internal/scaffold"
	"github.com/syed1006/goforge/internal/template"
	"github.com/syed1006/goforge/internal/templates"
)

type newOpts struct {
	seed          config.Config
	noInteractive bool
	dryRun        bool
	overwrite     bool
	latest        bool
}

func newNewCmd() *cobra.Command {
	o := &newOpts{seed: config.Config{Framework: "", Database: ""}}

	cmd := &cobra.Command{
		Use:   "new [project]",
		Short: "Create a new Go project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				o.seed.ProjectName = args[0]
			}
			return runNew(cmd, o)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&o.seed.ProjectName, "name", "", "Project name (overrides positional)")
	flags.StringVar(&o.seed.ModulePath, "module", "", "Go module path (e.g. github.com/me/myapi)")
	flags.StringVar(&o.seed.OutputDir, "out", "", "Output directory (defaults to project name)")
	flags.StringVar(&o.seed.GoVersion, "go", "", "Go directive in generated go.mod (default 1.23)")
	flags.StringVar((*string)(&o.seed.Framework), "framework", "", "HTTP framework: stdlib|chi|gin|fiber|echo")
	flags.StringVar((*string)(&o.seed.Database), "database", "", "Database driver: none|postgres|mysql|sqlite|mongo|redis")
	flags.BoolVar(&o.seed.GRPC, "grpc", false, "Add a gRPC server")
	flags.BoolVar(&o.seed.GraphQL, "graphql", false, "Add a GraphQL server (gqlgen)")
	flags.BoolVar(&o.seed.HotReload, "hot-reload", false, "Add hot-reload via air")
	flags.BoolVar(&o.seed.Lint, "lint", false, "Add golangci-lint and pre-commit")
	flags.BoolVar(&o.seed.Docker, "docker", false, "Add Dockerfile and docker-compose.yml")
	flags.BoolVar(&o.seed.CI, "ci", false, "Add GitHub Actions workflow")
	flags.BoolVar(&o.seed.OTel, "otel", false, "Add OpenTelemetry tracing (OTLP exporter)")
	flags.BoolVar(&o.seed.Metrics, "metrics", false, "Add Prometheus /metrics endpoint")
	flags.BoolVar(&o.noInteractive, "no-interactive", false, "Disable interactive prompts; all values must come from flags")
	flags.BoolVar(&o.dryRun, "dry-run", false, "Walk every generator but write nothing to disk")
	flags.BoolVar(&o.overwrite, "overwrite", false, "Allow writing into a non-empty output directory")
	flags.BoolVar(&o.latest, "latest", false, "Resolve every module dependency at @latest instead of the pinned versions")

	return cmd
}

func runNew(cmd *cobra.Command, o *newOpts) error {
	cfg := o.seed

	if !o.noInteractive {
		out, err := prompt.Ask(cfg)
		if err != nil {
			return err
		}
		cfg = out
	} else {
		if cfg.GoVersion == "" {
			cfg.GoVersion = "1.23"
		}
		if cfg.ProjectName == "" {
			return errors.New("--no-interactive requires --name or a positional project name")
		}
		if cfg.ModulePath == "" {
			return errors.New("--no-interactive requires --module")
		}
		if cfg.Framework == "" {
			cfg.Framework = config.FrameworkStdlib
		}
		if cfg.Database == "" {
			cfg.Database = config.DatabaseNone
		}
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), prompt.Summary(cfg))

	tplFS := templates.FS()
	eng, err := template.New(tplFS)
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	return scaffold.Run(cfg, generators.Default(), eng, scaffold.Options{
		DryRun:    o.dryRun,
		Overwrite: o.overwrite,
		Latest:    o.latest,
		Log:       os.Stdout,
	})
}
