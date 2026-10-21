// Package project detects btd projects: directories that contain a
// .btd/manifest.json marker.
package project

import (
	"os"
	"path/filepath"
)

// MarkerDir is the project-local directory that holds the project manifest
// and marks a directory as a btd project.
const MarkerDir = ".btd"

// MarkerFile is the manifest filename that identifies a btd project.
const MarkerFile = "manifest.json"

// Find walks up from start (empty means the current working directory) to the
// nearest ancestor containing a .btd/manifest.json. It returns that project
// root and true, or "" and false if no marker is found.
func Find(start string) (root string, ok bool) {
	if start == "" {
		start, _ = os.Getwd()
	}
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, MarkerDir, MarkerFile)); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
