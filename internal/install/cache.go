package install

import (
	"os"
	"path/filepath"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

// CacheDir returns the tomo clone cache directory, honoring XDG_CACHE_HOME.
func CacheDir() string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "btd", "tomo")
}

// gitStrategy installs a git-kind tool by cloning (or refreshing) the repo at
// its pinned ref into a managed cache.
type gitStrategy struct{}

func (gitStrategy) Install(r Runner, tool manifest.Tool, o Options) error {
	dest := o.CacheDir
	if dest == "" {
		dest = CacheDir()
	}
	return ensureClone(r, tool, dest)
}

// ensureClone clones tool.Repo into dest at tool.Ref, or fetches and checks
// out the ref if the checkout already exists.
func ensureClone(r Runner, tool manifest.Tool, dest string) error {
	if _, err := os.Stat(filepath.Join(dest, ".git")); err == nil {
		if err := r.Run(dest, "git", "fetch", "--depth", "1", "origin", tool.Ref); err != nil {
			return err
		}
		return r.Run(dest, "git", "checkout", "--detach", tool.Ref)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return r.Run("", "git", "clone", "--depth", "1", "--branch", tool.Ref, tool.Repo, dest)
}
