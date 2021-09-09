package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildScaffoldCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scaffold",
		Short: ScaffoldCommandsShortDescription,
	}

	cmd.AddCommand(buildScaffoldClusterCommand(o))
	cmd.AddCommand(buildScaffoldApplicationCommand(o))

	return cmd
}
