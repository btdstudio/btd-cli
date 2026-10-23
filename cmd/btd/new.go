package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
	newpkg "github.com/btdstudio/btd-cli/internal/new"
)

var (
	newInclude  string
	newSiteDir  string
	newProvider string
)

var newCmd = &cobra.Command{
	Use:   "new <template> <name> [module-path]",
	Short: "Scaffold a new project",
	Long: `Scaffold a new project with goph and optionally include tomo and hachigo.

goph new is non-interactive and fails if the target directory already exists.
Inclusions install per-project, never globally.`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNew(cmd, args)
	},
}

func init() {
	newCmd.Flags().StringVar(&newInclude, "include", "", "comma-separated tools to include (tomo,hachigo)")
	newCmd.Flags().StringVar(&newSiteDir, "site-dir", "", "hachigo site directory (default \"site\")")
	newCmd.Flags().StringVar(&newProvider, "provider", "all", "tomo harness selection (all|opencode|claude|copilot|antigravity|codex)")
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	template := args[0]
	name := args[1]
	modulePath := ""
	if len(args) == 3 {
		modulePath = args[2]
	}
	if template != "goph" {
		return fmt.Errorf("unsupported template %q (v1 supports only \"goph\")", template)
	}

	includes, err := parseIncludes(newInclude)
	if err != nil {
		return err
	}

	m, err := manifest.Load("")
	if err != nil {
		return err
	}

	o := newpkg.Options{
		Template:   template,
		Name:       name,
		ModulePath: modulePath,
		Include:    includes,
		SiteDir:    newSiteDir,
		Provider:   newProvider,
	}
	if err := newpkg.Run(install.SystemRunner{}, m, o); err != nil {
		return err
	}

	printNextSteps(cmd, o)
	return nil
}

func parseIncludes(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
	var out []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		if p != "tomo" && p != "hachigo" {
			return nil, fmt.Errorf("unknown include %q (expected tomo or hachigo)", p)
		}
		out = append(out, p)
	}
	return out, nil
}

func printNextSteps(cmd *cobra.Command, o newpkg.Options) {
	out := cmd.OutOrStdout()
	site := o.SiteDir
	if site == "" {
		site = "site"
	}
	fmt.Fprintf(out, "\nCreated %s.\n", o.Name)
	fmt.Fprintf(out, "Next steps:\n")
	fmt.Fprintf(out, "  cd %s\n", o.Name)
	fmt.Fprintf(out, "  cp .env.example .env\n")
	if has(o.Include, "hachigo") {
		fmt.Fprintf(out, "  hachigo serve -c %s/config.yml\n", site)
	}
	if has(o.Include, "tomo") {
		fmt.Fprintf(out, "  btd doctor    # verify the toolchain\n")
	}
}

func has(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
