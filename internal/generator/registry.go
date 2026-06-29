package generator

import "github.com/syed1006/goforge/internal/config"

// Registry holds an ordered collection of generators.
type Registry struct {
	gens []Generator
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry { return &Registry{} }

// Register appends g to the registry. Order of registration is preserved and is
// the order in which generators run, so callers should register foundational
// generators (base, framework) before ones that build on top (lint, CI).
func (r *Registry) Register(g ...Generator) {
	r.gens = append(r.gens, g...)
}

// All returns every registered generator.
func (r *Registry) All() []Generator { return r.gens }

// Applicable returns the subset of generators that opt into the given config.
func (r *Registry) Applicable(cfg config.Config) []Generator {
	out := make([]Generator, 0, len(r.gens))
	for _, g := range r.gens {
		if g.Applies(cfg) {
			out = append(out, g)
		}
	}
	return out
}
