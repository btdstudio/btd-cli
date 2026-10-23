package new

import "github.com/btdstudio/btd-cli/internal/manifest"

// seedManifest writes the project manifest into the scaffolded project,
// recording the toolset pinned at scaffold time.
func seedManifest(name string, m *manifest.Manifest) error {
	return manifest.WriteFile(manifest.ProjectPath(name), m)
}
