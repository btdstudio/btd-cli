package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/btdstudio/btd-cli/internal/install"
	"github.com/btdstudio/btd-cli/internal/manifest"
	"github.com/btdstudio/btd-cli/internal/update"
)

var updateCheck bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the toolchain pins to the latest compatible versions",
	Long: `Align the toolchain pins with btd's curated known-good set for this
release, then re-install any changed tools. Pins move to refs validated
against the current btd release, never blindly to @latest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdate(cmd.OutOrStdout())
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "dry run: print what would change")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(out io.Writer) error {
	m, err := manifest.Load("")
	if err != nil {
		return err
	}
	changes, err := update.Run(install.SystemRunner{}, m, updateCheck)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		fmt.Fprintln(out, "All tools are up to date.")
		return nil
	}
	for _, c := range changes {
		fmt.Fprintf(out, "%s: %s -> %s\n", c.Tool, c.From, c.To)
	}
	if updateCheck {
		fmt.Fprintln(out, "(dry run — no changes applied)")
	}
	return nil
}
