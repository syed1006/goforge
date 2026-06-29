// Package scaffold orchestrates generators to produce a complete project tree.
package scaffold

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/fsutil"
	"github.com/syed1006/goforge/internal/generator"
)

// Options tweak orchestrator behavior. The zero value is a sensible default.
type Options struct {
	// DryRun prevents disk mutation; the orchestrator still walks every
	// generator and renders every file, just doesn't commit anything.
	DryRun bool
	// Overwrite controls collision behavior when files already exist.
	Overwrite bool
	// SkipPostSteps disables go-mod-tidy and gofmt invocations after the writes.
	SkipPostSteps bool
	// Log is the destination for progress messages. Defaults to os.Stdout.
	Log io.Writer
}

// Run executes every applicable generator and writes the project tree.
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
	root, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("resolve output dir: %w", err)
	}
	if err := guardRoot(root, opts); err != nil {
		return err
	}

	mode := fsutil.ModeError
	if opts.Overwrite {
		mode = fsutil.ModeOverwrite
	}
	writer, err := fsutil.NewWriter(root,
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
	fmt.Fprintf(opts.Log, "→ scaffolding into %s\n", root)
	for _, g := range applicable {
		fmt.Fprintf(opts.Log, "→ %s\n", g.Name())
		if err := g.Generate(ctx); err != nil {
			return fmt.Errorf("%s: %w", g.Name(), err)
		}
	}

	if err := writeGoMod(writer, cfg, manifest); err != nil {
		return err
	}

	if opts.DryRun || opts.SkipPostSteps {
		return nil
	}
	return runPostSteps(root, manifest, opts.Log)
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

func runPostSteps(root string, m *generator.Manifest, log io.Writer) error {
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintln(log, "warning: go binary not found on PATH; skipping mod resolution")
		return nil
	}

	for _, req := range m.Requires() {
		spec := req.Module
		if req.Version != "" && req.Version != "latest" {
			spec = req.Module + "@" + req.Version
		} else {
			spec = req.Module + "@latest"
		}
		fmt.Fprintf(log, "→ go get %s\n", spec)
		cmd := exec.Command("go", "get", spec)
		cmd.Dir = root
		cmd.Stdout = log
		cmd.Stderr = log
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("go get %s: %w", spec, err)
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
