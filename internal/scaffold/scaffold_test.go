package scaffold

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/generator"
)

type fakeRenderer struct{}

func (fakeRenderer) Render(name string, _ any) ([]byte, error) {
	return []byte("// " + name), nil
}

type recordingGen struct {
	name      string
	apply     bool
	wroteFile string
	mods      [][2]string
	failWith  error
}

func (r *recordingGen) Name() string                 { return r.name }
func (r *recordingGen) Applies(_ config.Config) bool { return r.apply }
func (r *recordingGen) Generate(ctx *generator.Context) error {
	if r.wroteFile != "" {
		if err := ctx.Writer.Write(r.wroteFile, []byte("payload"), 0o644); err != nil {
			return err
		}
	}
	for _, m := range r.mods {
		ctx.Manifest.Require(m[0], m[1])
	}
	return r.failWith
}

func baseCfg(out string) config.Config {
	return config.Config{
		ProjectName: "myapi",
		ModulePath:  "github.com/me/myapi",
		GoVersion:   "1.23",
		Framework:   config.FrameworkStdlib,
		Database:    config.DatabaseNone,
		OutputDir:   out,
	}
}

func TestRunHappyPath(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "out")
	cfg := baseCfg(dir)

	reg := generator.NewRegistry()
	reg.Register(
		&recordingGen{name: "a", apply: true, wroteFile: "a.txt", mods: [][2]string{{"github.com/x/y", "v1.0.0"}}},
		&recordingGen{name: "skip", apply: false, wroteFile: "skip.txt"},
		&recordingGen{name: "b", apply: true, mods: [][2]string{{"github.com/p/q", "v0.1.0"}}},
	)

	var log bytes.Buffer
	if err := Run(cfg, reg, fakeRenderer{}, Options{Log: &log, SkipPostSteps: true}); err != nil {
		t.Fatalf("Run: %v\nlog: %s", err, log.String())
	}

	gomod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	for _, want := range []string{"module github.com/me/myapi", "go 1.23"} {
		if !bytes.Contains(gomod, []byte(want)) {
			t.Errorf("go.mod missing %q:\n%s", want, gomod)
		}
	}

	if _, err := os.Stat(filepath.Join(dir, "a.txt")); err != nil {
		t.Errorf("expected a.txt written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "skip.txt")); !os.IsNotExist(err) {
		t.Error("skip.txt should not exist")
	}
}

func TestRunRefusesNonEmptyDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "existing"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := baseCfg(dir)
	reg := generator.NewRegistry()
	if err := Run(cfg, reg, fakeRenderer{}, Options{SkipPostSteps: true, Log: io.Discard}); err == nil {
		t.Fatal("expected error for non-empty dir")
	}
}

func TestRunValidatesConfig(t *testing.T) {
	t.Parallel()
	cfg := baseCfg(t.TempDir())
	cfg.ProjectName = "Bad-Name"
	if err := Run(cfg, generator.NewRegistry(), fakeRenderer{}, Options{Log: io.Discard}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRunGeneratorFailureLeavesNoArtifacts(t *testing.T) {
	t.Parallel()
	parent := t.TempDir()
	dir := filepath.Join(parent, "out")
	cfg := baseCfg(dir)

	reg := generator.NewRegistry()
	reg.Register(
		&recordingGen{name: "ok", apply: true, wroteFile: "ok.txt"},
		&recordingGen{name: "boom", apply: true, failWith: errors.New("kaboom")},
	)

	err := Run(cfg, reg, fakeRenderer{}, Options{SkipPostSteps: true, Log: io.Discard})
	if err == nil {
		t.Fatal("expected error from failing generator")
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("final dir should not exist after failure, stat err: %v", err)
	}
	entries, err := os.ReadDir(parent)
	if err != nil {
		t.Fatalf("ReadDir parent: %v", err)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), "goforge-staging") {
			t.Errorf("staging dir leaked: %s", e.Name())
		}
	}
}

func TestRunPromotesIntoPreExistingEmptyDir(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "out")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := baseCfg(dir)
	reg := generator.NewRegistry()
	reg.Register(&recordingGen{name: "g", apply: true, wroteFile: "a.txt"})
	if err := Run(cfg, reg, fakeRenderer{}, Options{SkipPostSteps: true, Log: io.Discard}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "a.txt")); err != nil {
		t.Errorf("expected a.txt in final dir: %v", err)
	}
}
