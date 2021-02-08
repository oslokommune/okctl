package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildVersionCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version info",
		Long:  `Prints version number and git-sha of okctl-build.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(o.Out, version.String())
			return err
		},
	}

	return cmd
}
