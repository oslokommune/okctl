package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildApplyCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply commands",
	}

	cmd.AddCommand(buildApplyApplicationCommand(o))
	cmd.AddCommand(buildApplyClusterCommand(o))

	return cmd
}
