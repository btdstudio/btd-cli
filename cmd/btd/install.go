package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
	"github.com/btdstudio/btd-cli/internal/project"
)

var (
	installProvider string
	installYes      bool
	installSiteDir  string
)

var installCmd = &cobra.Command{
	Use:   "install <goph|tomo|hachigo|all>",
	Short: "Install a tool from the btdstudio toolchain",
	Long: `Install a tool from the btdstudio toolchain, pinned by the manifest.

Inside a btd project (.btd/manifest.json) installs are per-project and touch
only the project's own directories. Outside a project they install globally,
writing to user-level directories (GOPATH/bin, harness config), and prompt for
confirmation first (skip with --yes).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstall(cmd, args[0])
	},
}

func init() {
	installCmd.Flags().StringVar(&installProvider, "provider", "all", "tomo harness selection (all|opencode|claude|copilot|antigravity|codex)")
	installCmd.Flags().BoolVarP(&installYes, "yes", "y", false, "skip the global-install confirmation prompt")
	installCmd.Flags().StringVar(&installSiteDir, "site-dir", "", "hachigo site directory (default \"site\")")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, tool string) error {
	root, isProject := project.Find("")

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	o := install.Options{
		Provider:    installProvider,
		SiteDir:     installSiteDir,
		SkipConfirm: installYes,
	}
	if isProject {
		o.ProjectRoot = root
	}

	if !isProject && !installYes {
		names := []string{tool}
		if tool == "all" {
			names = []string{"goph", "hachigo", "tomo"}
		}
		if !install.Confirm(cmd.InOrStdin(), cmd.OutOrStdout(), m, names) {
			fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
			return nil
		}
	}

	return install.Install(install.SystemRunner{}, m, tool, o)
}
