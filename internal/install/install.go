// Package install implements btd's install strategies: the go and git kinds,
// plus context-aware installers for goph, hachigo, and tomo.
package install

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/btdstudio/btd-cli/internal/manifest"
)

// Runner executes a single external command with streamed stdio. Tests
// substitute a fake to avoid touching the real system.
type Runner interface {
	Run(dir string, name string, args ...string) error
}

// SystemRunner runs commands via os/exec, streaming stdio to the terminal and
// propagating non-zero exit codes.
type SystemRunner struct{}

func (SystemRunner) Run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Options carries the context for an install operation.
type Options struct {
	ProjectRoot string // non-empty installs per-project rather than globally
	Provider    string // tomo harness selection ("" = tomo default)
	SiteDir     string // hachigo site directory ("" = "site")
	CacheDir    string // tomo clone cache ("" = XDG default)
	SkipConfirm bool   // --yes
}

// Strategy installs one tool from its manifest entry. goStrategy and
// gitStrategy implement the two install kinds; hachigoStrategy and tomoStrategy
// compose them with per-tool steps.
type Strategy interface {
	Install(r Runner, tool manifest.Tool, o Options) error
}

// registry maps tool names to their install strategies.
var registry = map[string]Strategy{
	"goph":    goStrategy{},
	"hachigo": hachigoStrategy{},
	"tomo":    tomoStrategy{},
}

// Install installs the named tool (or "all") using the effective manifest.
func Install(r Runner, m *manifest.Manifest, tool string, o Options) error {
	names := []string{tool}
	if tool == "all" {
		names = []string{"goph", "hachigo", "tomo"}
	}
	for _, name := range names {
		t, ok := m.Tools[name]
		if !ok {
			return fmt.Errorf("tool %q is not in the manifest", name)
		}
		s, ok := registry[name]
		if !ok {
			return fmt.Errorf("unknown tool %q (expected goph, tomo, hachigo, or all)", name)
		}
		if err := s.Install(r, t, o); err != nil {
			return fmt.Errorf("install %s: %w", name, err)
		}
	}
	return nil
}

// LookPath reports whether name is on PATH; overridable in tests.
var LookPath = exec.LookPath
