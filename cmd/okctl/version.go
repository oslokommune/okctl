package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nolint: globalvar
// The current git commit at build-time with goreleaser
var GitCommit string

// nolint: globalvar
// The current version of OkCLI
var OkctlVersion string

func buildVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version info",
		Long:  `Prints version number and git-sha of okctl-build.`,
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Printf("Version: %s\n", OkctlVersion)
			fmt.Printf("Commit : %s\n", GitCommit)
			return nil
		},
	}

	return cmd
}
