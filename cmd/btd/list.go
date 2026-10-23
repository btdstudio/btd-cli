package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed tools and their versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(out io.Writer) error {
	m, err := manifest.Load("")
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%-8s %-12s %-12s %s\n", "TOOL", "PIN", "INSTALLED", "STATUS")
	for _, name := range []string{"kumi", "hachigo", "tomo"} {
		tool := m.Tools[name]
		pin := pinOf(tool)
		installed := resolveInstalled(name, tool)
		status := "missing"
		if installed != "" {
			if installed == pin {
				status = "ok"
			} else {
				status = "mismatch"
			}
		}
		fmt.Fprintf(out, "%-8s %-12s %-12s %s\n", name, pin, installed, status)
	}
	return nil
}

func pinOf(t manifest.Tool) string {
	if t.Kind == manifest.KindGit {
		return t.Ref
	}
	return t.Version
}

// resolveInstalled returns the installed version of a tool ("" if missing).
// Overridable in tests.
var resolveInstalled = func(name string, tool manifest.Tool) string {
	switch tool.Kind {
	case manifest.KindGo:
		return binaryVersion(name)
	case manifest.KindGit:
		return checkoutRef(install.CacheDir())
	default:
		return ""
	}
}

// versionRe extracts a semver-like token ("v0.1.0", or a pseudo-version like
// "v0.0.0-20261018…-sha+dirty") from a `--version` line such as
// "kumi version v0.1.0".
var versionRe = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+[^\s]*`)

// extractVersion returns the version token from a version string, or "".
func extractVersion(s string) string {
	return versionRe.FindString(s)
}

func binaryVersion(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return ""
	}
	return extractVersion(string(out))
}

func checkoutRef(cache string) string {
	if _, err := os.Stat(filepath.Join(cache, ".git")); err != nil {
		return ""
	}
	out, err := exec.Command("git", "-C", cache, "describe", "--tags", "--always").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
