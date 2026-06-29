package fsutil

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriterCreatesNested(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	w, err := NewWriter(dir)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if err := w.Write("a/b/c.txt", []byte("hi"), 0); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "a/b/c.txt"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "hi" {
		t.Fatalf("content: %q", got)
	}
}

func TestWriterCollisionError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	w, _ := NewWriter(dir, WithMode(ModeError))
	_ = w.Write("x.txt", []byte("a"), 0)
	if err := w.Write("x.txt", []byte("b"), 0); err == nil {
		t.Fatal("expected collision error")
	}
}

func TestWriterCollisionSkip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	w, _ := NewWriter(dir, WithMode(ModeSkip))
	_ = w.Write("x.txt", []byte("first"), 0)
	if err := w.Write("x.txt", []byte("second"), 0); err != nil {
		t.Fatalf("expected skip not to error: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "x.txt"))
	if string(got) != "first" {
		t.Fatalf("skip should preserve first write, got %q", got)
	}
}

func TestWriterDryRun(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	var log bytes.Buffer
	w, _ := NewWriter(dir, DryRun(true), WithLog(&log))
	if err := w.Write("x.txt", []byte("nope"), 0); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "x.txt")); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote file: %v", err)
	}
	if log.Len() == 0 {
		t.Fatal("dry run should still log")
	}
}
