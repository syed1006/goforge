// Package scaffold orchestrates generators to produce a complete project tree.
package scaffold

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/fsutil"
	"github.com/syed1006/goforge/internal/generator"
)

// Options tweak orchestrator behavior. The zero value is a sensible default.
type Options struct {
	DryRun        bool
	Overwrite     bool
	SkipPostSteps bool
	// Latest forces every manifest requirement to be fetched at @latest instead
	// of its pinned version. Use only when you accept the reproducibility cost.
	Latest bool
	Log    io.Writer
}

// Run executes every applicable generator and writes the project tree.
//
// Writes go to a sibling staging directory first; on success the staging
// directory is renamed onto the final path (atomic on the same filesystem).
// Any failure removes the staging directory, so the final path is either the
// fully-built project or untouched — never a partial scaffold.
func Run(cfg config.Config, reg *generator.Registry, renderer generator.Renderer, opts Options) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if opts.Log == nil {
		opts.Log = os.Stdout
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = cfg.ProjectName
	}
	final, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("resolve output dir: %w", err)
	}
	if err := guardRoot(final, opts); err != nil {
		return err
	}

	// Dry-run never touches disk, so we can let the writer pretend its root is
	// the final path. Real runs route every write through a staging dir.
	workDir := final
	promoted := false
	if !opts.DryRun {
		staging, err := makeStaging(final)
		if err != nil {
			return err
		}
		workDir = staging
		defer func() {
			if !promoted {
				_ = os.RemoveAll(staging)
			}
		}()
	}

	mode := fsutil.ModeError
	if opts.Overwrite {
		mode = fsutil.ModeOverwrite
	}
	writer, err := fsutil.NewWriter(workDir,
		fsutil.WithMode(mode),
		fsutil.WithLog(opts.Log),
		fsutil.DryRun(opts.DryRun),
	)
	if err != nil {
		return err
	}

	manifest := generator.NewManifest(cfg.GoVersion)
	ctx := &generator.Context{
		Config:   cfg,
		Renderer: renderer,
		Writer:   writer,
		Manifest: manifest,
	}

	applicable := reg.Applicable(cfg)
	fmt.Fprintf(opts.Log, "→ scaffolding into %s\n", final)
	for _, g := range applicable {
		fmt.Fprintf(opts.Log, "→ %s\n", g.Name())
		if err := g.Generate(ctx); err != nil {
			return fmt.Errorf("%s: %w", g.Name(), err)
		}
	}

	if err := writeGoMod(writer, cfg, manifest); err != nil {
		return err
	}

	if opts.DryRun {
		return nil
	}
	if !opts.SkipPostSteps {
		if err := runPostSteps(workDir, cfg, manifest, opts); err != nil {
			return err
		}
	}

	if err := promote(workDir, final, opts.Overwrite); err != nil {
		return err
	}
	promoted = true
	return nil
}

func makeStaging(final string) (string, error) {
	parent := filepath.Dir(final)
	base := filepath.Base(final)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", fmt.Errorf("create parent for staging: %w", err)
	}
	name := "." + base + ".goforge-staging-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	staging := filepath.Join(parent, name)
	if err := os.Mkdir(staging, 0o755); err != nil {
		return "", fmt.Errorf("create staging dir: %w", err)
	}
	return staging, nil
}

func promote(staging, final string, overwrite bool) error {
	info, err := os.Stat(final)
	switch {
	case errors.Is(err, os.ErrNotExist):
		// fall through to rename
	case err != nil:
		return err
	case !info.IsDir():
		return fmt.Errorf("%s exists and is not a directory", final)
	default:
		entries, derr := os.ReadDir(final)
		if derr != nil {
			return derr
		}
		switch {
		case len(entries) == 0:
			if rerr := os.Remove(final); rerr != nil {
				return fmt.Errorf("remove empty final dir: %w", rerr)
			}
		case overwrite:
			if rerr := os.RemoveAll(final); rerr != nil {
				return fmt.Errorf("remove non-empty final dir: %w", rerr)
			}
		default:
			return fmt.Errorf("output directory %s became non-empty during scaffold", final)
		}
	}
	return os.Rename(staging, final)
}

func guardRoot(root string, opts Options) error {
	info, err := os.Stat(root)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil
	case err != nil:
		return err
	case !info.IsDir():
		return fmt.Errorf("%s exists and is not a directory", root)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	if len(entries) > 0 && !opts.Overwrite && !opts.DryRun {
		return fmt.Errorf("output directory %s is not empty (pass --overwrite to allow)", root)
	}
	return nil
}

func writeGoMod(w *fsutil.Writer, cfg config.Config, m *generator.Manifest) error {
	var b strings.Builder
	fmt.Fprintf(&b, "module %s\n\ngo %s\n", cfg.ModulePath, m.GoVersion)
	return w.Write("go.mod", []byte(b.String()), 0o644)
}

func runPostSteps(root string, cfg config.Config, m *generator.Manifest, opts Options) error {
	log := opts.Log
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintln(log, "warning: go binary not found on PATH; skipping mod resolution")
		return nil
	}

	reqs := m.Requires()
	if len(reqs) > 0 {
		specs := make([]string, 0, len(reqs))
		for _, req := range reqs {
			version := req.Version
			if opts.Latest || version == "" || version == "latest" {
				version = "latest"
			}
			specs = append(specs, req.Module+"@"+version)
		}
		fmt.Fprintf(log, "→ go get %s\n", strings.Join(specs, " "))
		args := append([]string{"get"}, specs...)
		cmd := exec.Command("go", args...)
		cmd.Dir = root
		cmd.Stdout = log
		cmd.Stderr = log
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("go get: %w", err)
		}
	}

	if cfg.GraphQL {
		fmt.Fprintln(log, "→ gqlgen generate")
		cmd := exec.Command("go", "run", "github.com/99designs/gqlgen", "generate")
		cmd.Dir = root
		cmd.Stdout = log
		cmd.Stderr = log
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("gqlgen generate: %w", err)
		}
	}

	fmt.Fprintln(log, "→ go mod tidy")
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = root
	tidy.Stdout = log
	tidy.Stderr = log
	if err := tidy.Run(); err != nil {
		fmt.Fprintf(log, "warning: go mod tidy failed: %v\n", err)
	}

	if _, err := exec.LookPath("gofmt"); err == nil {
		cmd := exec.Command("gofmt", "-s", "-w", ".")
		cmd.Dir = root
		cmd.Stdout = log
		cmd.Stderr = log
		_ = cmd.Run()
	}
	return nil
}
