package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildForwardCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forward",
		Short: "Forward commands",
		Long: `Forward traffic to different services running in AWS

Sometimes it can be useful to forward traffic to services running
in AWS, such as a Postgres database, to inspect the schema 
or perform other simple operations.
`,
	}

	cmd.AddCommand(buildForwardPostgres(o))

	return cmd
}
