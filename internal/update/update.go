// Package update implements `btd update`: align the manifest pins with btd's
// curated known-good set and re-install changed tools.
package update

import (
	"time"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
)

// ToolOrder is the canonical order tools are reported in.
var ToolOrder = []string{"goph", "hachigo", "tomo"}

// Change is one tool whose pin would move.
type Change struct {
	Tool string
	From string
	To   string
}

// Plan returns the pins that differ between m and btd's curated defaults.
func Plan(m *manifest.Manifest) []Change {
	defaults := manifest.Default().Tools
	var changes []Change
	for _, name := range ToolOrder {
		cur, ok := m.Tools[name]
		if !ok {
			continue
		}
		if from, to := pin(cur), pin(defaults[name]); from != to {
			changes = append(changes, Change{Tool: name, From: from, To: to})
		}
	}
	return changes
}

// Run aligns the manifest pins with the curated defaults. With check it only
// reports; otherwise it writes the new pins to the global manifest and
// re-installs changed tools.
func Run(r install.Runner, m *manifest.Manifest, check bool) ([]Change, error) {
	changes := Plan(m)
	if check || len(changes) == 0 {
		return changes, nil
	}
	if err := writeGlobalPins(changes); err != nil {
		return nil, err
	}
	defaults := manifest.Default()
	for _, c := range changes {
		if err := install.Install(r, defaults, c.Tool, install.Options{}); err != nil {
			return nil, err
		}
	}
	return changes, nil
}

// writeGlobalPins moves the changed tools' pins (plus the compatible label
// and updated timestamp) into the global manifest.
func writeGlobalPins(changes []Change) error {
	global, err := manifest.ReadFile(manifest.GlobalPath())
	if err != nil {
		return err
	}
	if global.Tools == nil {
		global.Tools = map[string]manifest.Tool{}
	}
	defaults := manifest.Default()
	for _, c := range changes {
		global.Tools[c.Tool] = defaults.Tools[c.Tool]
	}
	global.SchemaVersion = defaults.SchemaVersion
	global.Compatible = defaults.Compatible
	global.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return manifest.WriteFile(manifest.GlobalPath(), global)
}

func pin(t manifest.Tool) string {
	if t.Kind == manifest.KindGit {
		return t.Ref
	}
	return t.Version
}
