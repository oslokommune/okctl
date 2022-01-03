package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildMaintenanceStateUploadCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state-upload <path to cluster.yaml>",
		Short: "uploads a state.db",
		Long:  longMaintenanceStateUploadDescription,
		Args:  cobra.ExactArgs(1),
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionMaintenanceStateUpload),
			hooks.InitializeOkctl(o),
		),
		RunE: func(_ *cobra.Command, args []string) error {
			path := args[0]

			fmt.Fprintln(o.Out, "Uploading state")

			f, err := o.FileSystem.Open(path)
			if err != nil {
				return fmt.Errorf("opening file for reading: %w", err)
			}

			defer func() {
				_ = f.Close()
			}()

			services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
			if err != nil {
				return fmt.Errorf("acquiring client services: %w", err)
			}

			err = services.RemoteState.Upload(api.ID{
				Region:       o.Declaration.Metadata.Region,
				AWSAccountID: o.Declaration.Metadata.AccountID,
				ClusterName:  o.Declaration.Metadata.Name,
			}, f)
			if err != nil {
				return fmt.Errorf("uploading state: %w", err)
			}

			fmt.Fprintln(o.Out, "successfully uploaded state")

			return nil
		},
		PostRunE: hooks.RunECombinator(
			hooks.EmitEndCommandExecutionEvent(metrics.ActionMaintenanceStateUpload),
		),
	}

	return cmd
}

const longMaintenanceStateUploadDescription = `
Uploads state to the remote state location.
`
