package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildForwardCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forward",
		Short: ForwardCommandsShortDescription,
		Long:  ForwardCommandsLongDescription,
	}

	cmd.AddCommand(buildForwardPostgres(o))
	addCommonCommandFlags(cmd)

	return cmd
}
