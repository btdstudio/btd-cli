package main

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version is the current version of the CLI, set at build time with ldflags
// (-X main.Version=...), or resolved at runtime from debug.ReadBuildInfo()
// when installed via go install.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "btd",
	Short: "Orchestrate the btdstudio toolchain (kumi, tomo, hachigo)",
	Long: `btd is a thin CLI that coordinates the btdstudio toolchain — kumi,
tomo, and hachigo — by wrapping their existing CLIs and install flows.

It holds no templates and no tool logic; every operation shells out to
the underlying tool.`,
}

func init() {
	resolveVersion()
	rootCmd.Version = Version
}

// resolveVersion fills Version when it is still the default "dev", preferring
// the module version from build info (go install) and falling back to the VCS
// revision with a dirty marker.
func resolveVersion() {
	if Version != "dev" {
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		Version = v
		return
	}
	var revision string
	var modified bool
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}
	if revision == "" {
		return
	}
	short := revision
	if len(short) > 7 {
		short = short[:7]
	}
	Version = "dev-" + short
	if modified {
		Version += "+dirty"
	}
}

// Execute runs the root command and returns any error it produces.
func Execute() error {
	return rootCmd.Execute()
}
