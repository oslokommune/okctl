package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildMaintenanceCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "maintenance",
		Short:  "contains various maintenance commands",
		Hidden: true,
	}

	cmd.AddCommand(buildMaintenanceStateAcquireLockCommand(o))
	cmd.AddCommand(buildMaintenanceStateReleaseLockCommand(o))
	cmd.AddCommand(buildMaintenanceStateDownloadCommand(o))
	cmd.AddCommand(buildMaintenanceStateUploadCommand(o))

	return cmd
}
