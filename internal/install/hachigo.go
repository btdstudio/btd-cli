package install

import (
	"path/filepath"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

// hachigoStrategy installs hachigo. Globally it is a plain go install; in a
// project it also scaffolds the site with `hachigo init <site-dir>`.
type hachigoStrategy struct{}

func (hachigoStrategy) Install(r Runner, tool manifest.Tool, o Options) error {
	if o.ProjectRoot == "" {
		return goStrategy{}.Install(r, tool, o)
	}
	if _, err := LookPath("hachigo"); err != nil {
		if err := (goStrategy{}).Install(r, tool, o); err != nil {
			return err
		}
	}
	siteDir := o.SiteDir
	if siteDir == "" {
		siteDir = "site"
	}
	return r.Run("", "hachigo", "init", filepath.Join(o.ProjectRoot, siteDir))
}
