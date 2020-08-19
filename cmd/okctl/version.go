package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// nolint: globalvar
// The current git commit at build-time with goreleaser
var GitCommit string

// nolint: globalvar
// The current version of okctl
var OkctlVersion string

func buildVersionCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version info",
		Long:  `Prints version number and git-sha of okctl-build.`,
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Fprintf(o.Out, "Version: %s\n", OkctlVersion)
			fmt.Fprintf(o.Out, "Commit : %s\n", GitCommit)
			return nil
		},
	}

	return cmd
}
