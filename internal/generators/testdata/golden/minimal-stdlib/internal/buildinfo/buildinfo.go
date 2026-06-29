// Package buildinfo exposes build-time metadata stamped via -ldflags.
package buildinfo

// Version is overridden at build time via `-X .../buildinfo.Version=$(VERSION)`.
var Version = "dev"
