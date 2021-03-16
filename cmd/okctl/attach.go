package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildAttachCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach commands",
		Long: `Attach to different services running in AWS

Sometimes it can be useful to attach to service running
in AWS, such as a Postgres database, to inspect the schema 
or perform other simple operations.
`,
	}

	cmd.AddCommand(buildAttachPostgres(o))

	return cmd
}
