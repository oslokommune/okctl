package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildAttachCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: AttachCommandsShortDescription,
		Long:  AttachCommandsLongDescription,
	}

	cmd.AddCommand(buildAttachPostgres(o))

	return cmd
}
