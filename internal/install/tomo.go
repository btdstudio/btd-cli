package install

import (
	"path/filepath"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

// tomoStrategy installs tomo (a Node source tool with a git-only
// distribution): clone/refresh the checkout into the cache, then run its
// install.js in project or global mode.
type tomoStrategy struct{}

func (tomoStrategy) Install(r Runner, tool manifest.Tool, o Options) error {
	dest := o.CacheDir
	if dest == "" {
		dest = CacheDir()
	}
	if err := ensureClone(r, tool, dest); err != nil {
		return err
	}

	args := []string{filepath.Join(dest, "scripts", "install.js")}
	if o.ProjectRoot != "" {
		args = append(args, "--project", o.ProjectRoot)
	} else {
		args = append(args, "--global")
	}
	if o.Provider != "" {
		args = append(args, "--provider", o.Provider)
	}
	return r.Run("", "node", args...)
}
