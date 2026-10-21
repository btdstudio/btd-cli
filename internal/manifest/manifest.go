// Package manifest models the btd manifest: the JSON file that pins the
// sibling tools (goph, tomo, hachigo) to mutually-compatible versions. A
// global manifest lives at ~/.config/btd/manifest.json; a project manifest at
// <root>/.btd/manifest.json. Effective settings merge defaults ← global ←
// project (project wins).
package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Install kinds select the strategy btd uses to install a tool.
const (
	KindGo  = "go"  // go install <path>@<version>
	KindGit = "git" // git clone/checkout <repo> at <ref>
)

// Manifest is the btd toolset pin. Tools is keyed by tool name ("goph",
// "tomo", "hachigo").
type Manifest struct {
	SchemaVersion int             `json:"schemaVersion"`
	Compatible    string          `json:"compatible"`
	UpdatedAt     string          `json:"updatedAt,omitempty"`
	Tools         map[string]Tool `json:"tools"`
}

// Tool is one pinned tool. Kind selects the install strategy: "go" uses
// Path+Version (go install), "git" uses Repo+Ref (clone/checkout).
type Tool struct {
	Kind    string `json:"kind"`
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Ref     string `json:"ref,omitempty"`
}

// Default returns the base manifest: schema v1, curated for the current btd
// release, with the three tools pinned to v0.1.0.
func Default() *Manifest {
	return &Manifest{
		SchemaVersion: 1,
		Compatible:    "btd-0.1.0",
		Tools: map[string]Tool{
			"goph":    {Kind: KindGo, Path: "github.com/btdstudio/goph/cmd/goph", Version: "v0.1.0"},
			"hachigo": {Kind: KindGo, Path: "github.com/btdstudio/hachigo/cmd/hachigo", Version: "v0.1.0"},
			"tomo":    {Kind: KindGit, Repo: "https://github.com/btdstudio/tomo.git", Ref: "v0.1.0"},
		},
	}
}

// GlobalPath returns the global manifest path, honoring XDG_CONFIG_HOME.
func GlobalPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		base = filepath.Join(userHome(), ".config")
	}
	return filepath.Join(base, "btd", "manifest.json")
}

// ProjectPath returns the project manifest path for a project root.
func ProjectPath(root string) string {
	return filepath.Join(root, ".btd", "manifest.json")
}

func userHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// ReadFile reads and parses a manifest. A missing file yields an empty
// manifest (not an error), so callers can layer it into a merge.
func ReadFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return &Manifest{}, nil
	}
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &m, nil
}

// WriteFile writes a manifest as indented JSON, creating parent directories.
func WriteFile(path string, m *Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads the global and project manifests and returns the effective
// manifest (defaults ← global ← project).
func Load(root string) (*Manifest, error) {
	global, err := ReadFile(GlobalPath())
	if err != nil {
		return nil, err
	}
	project, err := ReadFile(ProjectPath(root))
	if err != nil {
		return nil, err
	}
	return Merge(Default(), global, project), nil
}

// Merge overlays layers in order (later wins) and returns a new manifest.
// Tools merge per-tool and per-field, so a project can override one tool's
// version while inheriting the rest.
func Merge(layers ...*Manifest) *Manifest {
	out := &Manifest{Tools: map[string]Tool{}}
	for _, layer := range layers {
		if layer == nil {
			continue
		}
		if layer.SchemaVersion > 0 {
			out.SchemaVersion = layer.SchemaVersion
		}
		if layer.Compatible != "" {
			out.Compatible = layer.Compatible
		}
		if layer.UpdatedAt != "" {
			out.UpdatedAt = layer.UpdatedAt
		}
		for name, tool := range layer.Tools {
			merged := out.Tools[name]
			mergeTool(&merged, tool)
			out.Tools[name] = merged
		}
	}
	return out
}

func mergeTool(dst *Tool, src Tool) {
	if src.Kind != "" {
		dst.Kind = src.Kind
	}
	if src.Path != "" {
		dst.Path = src.Path
	}
	if src.Version != "" {
		dst.Version = src.Version
	}
	if src.Repo != "" {
		dst.Repo = src.Repo
	}
	if src.Ref != "" {
		dst.Ref = src.Ref
	}
}
