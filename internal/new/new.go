// Package new implements `btd new`: scaffold a project with kumi and
// optionally include tomo and hachigo.
package new

import (
	"fmt"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
)

// Options carries the parameters for `btd new`.
type Options struct {
	Template   string   // v1: "kumi"
	Name       string   // directory name (and default module path)
	ModulePath string   // optional module path override
	Include    []string // subset of {"tomo", "hachigo"}
	SiteDir    string   // hachigo site directory ("" = "site")
	Provider   string   // tomo harness selection ("" = "all")
	CacheDir   string   // tomo clone cache ("" = XDG default)
}

// Run scaffolds a new project: ensure kumi, run `kumi new`, seed the project
// manifest, then apply any --include tools.
func Run(r install.Runner, m *manifest.Manifest, o Options) error {
	if err := ensureGoph(r, m); err != nil {
		return err
	}

	args := []string{"new", o.Name}
	if o.ModulePath != "" {
		args = append(args, o.ModulePath)
	}
	if err := r.Run("", "kumi", args...); err != nil {
		return fmt.Errorf("kumi new: %w", err)
	}

	if err := seedManifest(o.Name, m); err != nil {
		return err
	}

	for _, inc := range o.Include {
		switch inc {
		case "hachigo":
			if err := includeHachigo(r, m, o); err != nil {
				return err
			}
		case "tomo":
			if err := includeTomo(r, m, o); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown include %q (expected tomo or hachigo)", inc)
		}
	}
	return nil
}

// ensureGoph installs kumi if it is not already on PATH.
func ensureGoph(r install.Runner, m *manifest.Manifest) error {
	if _, err := install.LookPath("kumi"); err == nil {
		return nil
	}
	return install.Install(r, m, "kumi", install.Options{})
}
