// Package template renders the embedded scaffold templates.
package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"text/template"
)

// Engine renders templates that have been pre-parsed from an io/fs.FS.
//
// Templates are referenced by their path relative to the root of the FS, with
// the `.tmpl` suffix stripped. So `base/main.go.tmpl` is rendered as `base/main.go`.
type Engine struct {
	tmpl *template.Template
}

// New parses every *.tmpl under root in the provided FS and returns an Engine.
func New(root fs.FS) (*Engine, error) {
	t := template.New("goforge").Funcs(FuncMap()).Option("missingkey=error")

	err := fs.WalkDir(root, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".tmpl") {
			return nil
		}
		b, rerr := fs.ReadFile(root, p)
		if rerr != nil {
			return fmt.Errorf("read template %q: %w", p, rerr)
		}
		name := strings.TrimSuffix(p, ".tmpl")
		if _, perr := t.New(name).Parse(string(b)); perr != nil {
			return fmt.Errorf("parse template %q: %w", p, perr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Engine{tmpl: t}, nil
}

// Render evaluates the named template with data and returns the resulting bytes.
func (e *Engine) Render(name string, data any) ([]byte, error) {
	name = path.Clean(name)
	t := e.tmpl.Lookup(name)
	if t == nil {
		return nil, fmt.Errorf("template %q not found", name)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template %q: %w", name, err)
	}
	return buf.Bytes(), nil
}

// Names returns every template name registered with the engine (useful for diagnostics).
func (e *Engine) Names() []string {
	var out []string
	for _, t := range e.tmpl.Templates() {
		if t.Name() == "goforge" {
			continue
		}
		out = append(out, t.Name())
	}
	return out
}
