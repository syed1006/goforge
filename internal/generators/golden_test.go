package generators

import (
	"flag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/syed1006/goforge/internal/config"
	"github.com/syed1006/goforge/internal/scaffold"
	"github.com/syed1006/goforge/internal/template"
	"github.com/syed1006/goforge/internal/templates"
)

var update = flag.Bool("update", false, "regenerate golden fixtures")

func TestGenerators_Golden(t *testing.T) {
	scenarios := []struct {
		name string
		cfg  config.Config
	}{
		{
			name: "minimal-stdlib",
			cfg: config.Config{
				ProjectName: "demo", ModulePath: "example.com/demo", GoVersion: "1.23",
				Framework: config.FrameworkStdlib, Database: config.DatabaseNone,
			},
		},
		{
			name: "gin-postgres-grpc",
			cfg: config.Config{
				ProjectName: "demo", ModulePath: "example.com/demo", GoVersion: "1.23",
				Framework: config.FrameworkGin, Database: config.DatabasePostgres,
				GRPC: true,
			},
		},
		{
			name: "chi-redis-lint-docker-ci",
			cfg: config.Config{
				ProjectName: "demo", ModulePath: "example.com/demo", GoVersion: "1.23",
				Framework: config.FrameworkChi, Database: config.DatabaseRedis,
				Lint: true, Docker: true, CI: true,
			},
		},
	}

	eng, err := template.New(templates.FS())
	if err != nil {
		t.Fatalf("template.New: %v", err)
	}

	for _, sc := range scenarios {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			t.Parallel()

			outDir := filepath.Join(t.TempDir(), "out")
			cfg := sc.cfg
			cfg.OutputDir = outDir

			if err := scaffold.Run(cfg, Default(), eng, scaffold.Options{
				SkipPostSteps: true,
				Log:           io.Discard,
			}); err != nil {
				t.Fatalf("scaffold.Run: %v", err)
			}

			actual := readTree(t, outDir)
			fixtureDir := filepath.Join("testdata", "golden", sc.name)

			if *update {
				if err := os.RemoveAll(fixtureDir); err != nil {
					t.Fatalf("RemoveAll fixture: %v", err)
				}
				writeTree(t, fixtureDir, actual)
				return
			}

			expected := readTree(t, fixtureDir)
			compareTrees(t, expected, actual)
		})
	}
}

func readTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := make(map[string]string)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		out[rel] = string(b)
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("walk %s: %v", root, err)
	}
	return out
}

func writeTree(t *testing.T, root string, tree map[string]string) {
	t.Helper()
	for path, body := range tree {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", full, err)
		}
	}
}

func compareTrees(t *testing.T, want, got map[string]string) {
	t.Helper()
	keys := func(m map[string]string) []string {
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return ks
	}

	wantKeys := keys(want)
	gotKeys := keys(got)

	for _, k := range wantKeys {
		if _, ok := got[k]; !ok {
			t.Errorf("missing file in scaffold output: %s", k)
		}
	}
	for _, k := range gotKeys {
		if _, ok := want[k]; !ok {
			t.Errorf("unexpected file in scaffold output: %s\nrun: go test -update", k)
		}
	}
	for _, k := range wantKeys {
		gv, ok := got[k]
		if !ok {
			continue
		}
		if want[k] != gv {
			t.Errorf("content mismatch for %s\nrun: go test -update\nfirst diff line:\n  want: %s\n  got:  %s",
				k, firstDiffLine(want[k], gv), firstDiffLine(gv, want[k]))
		}
	}
}

func firstDiffLine(a, b string) string {
	la := strings.Split(a, "\n")
	lb := strings.Split(b, "\n")
	for i := range la {
		if i >= len(lb) || la[i] != lb[i] {
			return la[i]
		}
	}
	return ""
}
