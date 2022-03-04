package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	deleteClusterArgs = 0
)

func buildDeleteCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: DeleteCommandsShortDescription,
	}

	cmd.AddCommand(buildDeleteClusterCommand(o))
	cmd.AddCommand(buildDeleteApplicationCommand(o))

	return cmd
}
