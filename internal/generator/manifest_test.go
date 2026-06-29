package generator

import "testing"

func TestManifestRequireSorts(t *testing.T) {
	t.Parallel()
	m := NewManifest("1.23")
	m.Require("github.com/zzz/x", "v1.0.0")
	m.Require("github.com/aaa/y", "v0.1.0")
	got := m.Requires()
	if len(got) != 2 || got[0].Module != "github.com/aaa/y" {
		t.Fatalf("Requires(): expected sorted output, got %+v", got)
	}
}

func TestManifestRequireOverride(t *testing.T) {
	t.Parallel()
	m := NewManifest("1.23")
	m.Require("github.com/x/y", "v1.0.0")
	m.Require("github.com/x/y", "v1.2.0")
	reqs := m.Requires()
	if len(reqs) != 1 || reqs[0].Version != "v1.2.0" {
		t.Fatalf("Require should override: got %+v", reqs)
	}
}

func TestRenderRequires(t *testing.T) {
	t.Parallel()
	m := NewManifest("1.23")
	if got := m.RenderRequires(); got != "" {
		t.Fatalf("empty manifest should render empty, got %q", got)
	}
	m.Require("github.com/a/b", "v1.0.0")
	want := "require (\n\tgithub.com/a/b v1.0.0\n)\n"
	if got := m.RenderRequires(); got != want {
		t.Fatalf("RenderRequires:\n got: %q\nwant: %q", got, want)
	}
}
