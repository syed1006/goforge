// Package fsutil contains the file-system writer used by the scaffold orchestrator.
package fsutil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Mode controls how the writer treats existing files.
type Mode int

const (
	// ModeError aborts when a file already exists at the destination path.
	ModeError Mode = iota
	// ModeOverwrite replaces any existing file at the destination path.
	ModeOverwrite
	// ModeSkip leaves any existing file at the destination path unchanged.
	ModeSkip
)

// Writer commits files relative to a fixed root directory.
type Writer struct {
	root     string
	mode     Mode
	logSink  io.Writer
	dryRun   bool
}

// Option configures a Writer.
type Option func(*Writer)

// WithMode controls the collision behavior of subsequent Write calls.
func WithMode(m Mode) Option { return func(w *Writer) { w.mode = m } }

// WithLog routes per-write activity messages to the provided sink.
func WithLog(sink io.Writer) Option { return func(w *Writer) { w.logSink = sink } }

// DryRun prevents any disk mutation but still logs the intended writes.
func DryRun(on bool) Option { return func(w *Writer) { w.dryRun = on } }

// NewWriter returns a Writer rooted at root. The directory is created if absent.
func NewWriter(root string, opts ...Option) (*Writer, error) {
	w := &Writer{root: root}
	for _, opt := range opts {
		opt(w)
	}
	if !w.dryRun {
		if err := os.MkdirAll(root, 0o755); err != nil {
			return nil, fmt.Errorf("create root %q: %w", root, err)
		}
	}
	return w, nil
}

// Root returns the absolute (or unmodified) project root.
func (w *Writer) Root() string { return w.root }

// Write commits content to relPath, honoring the configured collision mode.
func (w *Writer) Write(relPath string, content []byte, mode fs.FileMode) error {
	if relPath == "" {
		return errors.New("empty path")
	}
	target := filepath.Join(w.root, filepath.FromSlash(relPath))
	if mode == 0 {
		mode = 0o644
	}

	if exists(target) {
		switch w.mode {
		case ModeError:
			return fmt.Errorf("refusing to overwrite %q", target)
		case ModeSkip:
			w.log("skip   %s\n", relPath)
			return nil
		}
	}

	w.log("write  %s\n", relPath)
	if w.dryRun {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("mkdir for %q: %w", target, err)
	}
	if err := os.WriteFile(target, content, mode); err != nil {
		return fmt.Errorf("write %q: %w", target, err)
	}
	return nil
}

func (w *Writer) log(format string, args ...any) {
	if w.logSink == nil {
		return
	}
	fmt.Fprintf(w.logSink, format, args...)
}

func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
