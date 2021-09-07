package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildCreateCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: CreateCommandsShortDescription,
		Long:  CreateCommandsLongDescription,
	}

	cmd.AddCommand(buildCreatePostgresCommand(o))

	return cmd
}
