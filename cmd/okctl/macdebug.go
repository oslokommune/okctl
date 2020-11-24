package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

//
//// nolint: globalvar
//// The current git commit at build-time with goreleaser
//var GitCommit string
//
//// nolint: globalvar
//// The current version of okctl
//var OkctlVersion string

func buildMacDebugCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "macdebug",
		Short: "Check if this functionality exists",
		Long:  `I want to see whether or not the release creates this functionality as part of the binary.`,
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Fprintf(o.Out, "This functionality actually exists!")
			return nil
		},
	}

	return cmd
}
