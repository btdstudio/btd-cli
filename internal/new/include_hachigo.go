package new

import (
	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
)

// includeHachigo scaffolds a hachigo site into the project via the project-
// scoped install flow (ensure binary, then `hachigo init <site-dir>`).
func includeHachigo(r install.Runner, m *manifest.Manifest, o Options) error {
	return install.Install(r, m, "hachigo", install.Options{
		ProjectRoot: o.Name,
		SiteDir:     o.SiteDir,
	})
}
