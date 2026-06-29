package generator

import (
	"fmt"
	"sort"
	"strings"
)

// Manifest collects go.mod requirements declared by generators.
type Manifest struct {
	GoVersion string
	requires  map[string]string
	tools     map[string]string
}

// NewManifest returns a manifest pinned to the given go directive.
func NewManifest(goVersion string) *Manifest {
	return &Manifest{
		GoVersion: goVersion,
		requires:  make(map[string]string),
		tools:     make(map[string]string),
	}
}

// Require records a runtime dependency; later calls override earlier ones.
func (m *Manifest) Require(module, version string) {
	m.requires[module] = version
}

// Tool records a developer-tool install hint (not a runtime dependency).
func (m *Manifest) Tool(name, hint string) {
	m.tools[name] = hint
}

// Requires returns the runtime requirements sorted by module path.
func (m *Manifest) Requires() []ModuleRequirement {
	out := make([]ModuleRequirement, 0, len(m.requires))
	for k, v := range m.requires {
		out = append(out, ModuleRequirement{Module: k, Version: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Module < out[j].Module })
	return out
}

// Tools returns developer tool hints sorted by name.
func (m *Manifest) Tools() []ToolHint {
	out := make([]ToolHint, 0, len(m.tools))
	for k, v := range m.tools {
		out = append(out, ToolHint{Name: k, Hint: v})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ModuleRequirement is a single line of the go.mod require block.
type ModuleRequirement struct {
	Module  string
	Version string
}

// ToolHint is an installation hint for a developer tool.
type ToolHint struct {
	Name string
	Hint string
}

// String renders the require block as it should appear in go.mod.
func (m *Manifest) RenderRequires() string {
	reqs := m.Requires()
	if len(reqs) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("require (\n")
	for _, r := range reqs {
		fmt.Fprintf(&b, "\t%s %s\n", r.Module, r.Version)
	}
	b.WriteString(")\n")
	return b.String()
}
