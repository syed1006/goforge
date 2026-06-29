// Package templates exposes the embedded scaffold templates as an io/fs.FS.
package templates

import (
	"embed"
	"io/fs"
)

//go:embed all:files
var raw embed.FS

// FS returns an fs.FS rooted at the directory that holds template files.
func FS() fs.FS {
	sub, err := fs.Sub(raw, "files")
	if err != nil {
		panic(err)
	}
	return sub
}
