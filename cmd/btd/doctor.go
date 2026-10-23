package main

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
)

var doctorInstall bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the btdstudio toolchain for missing or unhealthy tools",
	Long: `Check that btd and the tools it coordinates (kumi, tomo, hachigo) and
their prerequisites are installed. With --install, btd installs the tools it
can and delegates the rest to kumi doctor --install.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDoctor(cmd)
	},
}

func init() {
	doctorCmd.Flags().BoolVarP(&doctorInstall, "install", "i", false, "install missing tools")
	rootCmd.AddCommand(doctorCmd)
}

type toolCheck struct {
	name    string
	ok      bool
	version string
}

func runDoctor(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()

	checks := []toolCheck{
		checkBinary("go"),
		checkBinary("git"),
		checkNode(),
		checkBinary("kumi"),
		checkBinary("hachigo"),
		checkBinary("templ"),
		checkBinary("air"),
		checkBinary("goose"),
	}
	for _, c := range checks {
		status := "ok"
		if !c.ok {
			status = "MISSING"
		}
		fmt.Fprintf(out, "%-8s %-8s %s\n", c.name, status, c.version)
	}

	if !doctorInstall {
		return nil
	}

	// kumi/hachigo install via the manifest; templ/air/goose delegate to
	// kumi doctor --install. go/git/node are system prerequisites btd cannot
	// install, so they are reported but left to the user.
	if !checks[3].ok {
		if err := installTool(out, "kumi"); err != nil {
			return err
		}
	}
	if !checks[4].ok {
		if err := installTool(out, "hachigo"); err != nil {
			return err
		}
	}
	if !checks[5].ok || !checks[6].ok || !checks[7].ok {
		if err := exec.Command("kumi", "doctor", "--install").Run(); err != nil {
			return fmt.Errorf("kumi doctor --install: %w", err)
		}
	}
	return nil
}

func installTool(out io.Writer, name string) error {
	fmt.Fprintf(out, "Installing %s...\n", name)
	m, err := manifest.Load("")
	if err != nil {
		return err
	}
	return install.Install(install.SystemRunner{}, m, name, install.Options{})
}

func checkBinary(name string) toolCheck {
	if _, err := exec.LookPath(name); err != nil {
		return toolCheck{name: name}
	}
	return toolCheck{name: name, ok: true, version: firstLine(runVersion(name))}
}

func checkNode() toolCheck {
	c := checkBinary("node")
	if !c.ok {
		return c
	}
	if major := nodeMajor(c.version); major < 20 {
		c.ok = false
		c.version = fmt.Sprintf("%s (want >=20, got %d)", c.version, major)
	}
	return c
}

func runVersion(name string) string {
	out, err := exec.Command(name, "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func nodeMajor(v string) int {
	// "v20.11.0" -> 20
	s := strings.TrimPrefix(v, "v")
	if i := strings.IndexByte(s, '.'); i >= 0 {
		s = s[:i]
	}
	n, _ := strconv.Atoi(s)
	return n
}
