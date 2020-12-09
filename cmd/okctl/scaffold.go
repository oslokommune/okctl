package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildScaffoldCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scaffold",
		Short: "scaffold templates for different resources",
	}

	cmd.AddCommand(buildScaffoldClusterCommand(o))

	return cmd
}
