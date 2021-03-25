package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildCreateCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create commands",
		Long: `Create various kinds of resources.

Pre-requisites:
okctl creates various configuration files, and assumes that it is
running in a git repository. Initialize or clone a git repository
before running any of these commands.`,
	}

	cmd.AddCommand(buildCreatePostgresCommand(o))

	return cmd
}
