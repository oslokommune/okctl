package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type stateDownloadOpts struct {
	ClusterDeclarationPath string
}

// nolint: funlen
func buildMaintenanceStateDownloadCommand(o *okctl.Okctl) *cobra.Command {
	opts := &stateDownloadOpts{}

	var path string

	cmd := &cobra.Command{
		Use:   "state-download",
		Short: "downloads a state.db",
		Long:  longMaintenanceStateDownloadDescription,
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionMaintenanceStateDownload),
			hooks.LoadClusterDeclaration(o, &opts.ClusterDeclarationPath),
			hooks.InitializeOkctl(o),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Fprintln(o.Out, "Downloading state")

			services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
			if err != nil {
				return fmt.Errorf("acquiring client services: %w", err)
			}

			state, err := services.RemoteState.Download(api.ID{
				Region:       o.Declaration.Metadata.Region,
				AWSAccountID: o.Declaration.Metadata.AccountID,
				ClusterName:  o.Declaration.Metadata.Name,
			})
			if err != nil {
				return fmt.Errorf("downloading state: %w", err)
			}

			err = o.FileSystem.WriteReader(path, state)
			if err != nil {
				return fmt.Errorf("writing state to file: %w", err)
			}

			fmt.Fprintf(o.Out, "successfully downloaded state to %s\n", path)

			return nil
		},
		PostRunE: hooks.RunECombinator(
			hooks.EmitEndCommandExecutionEvent(metrics.ActionMaintenanceStateDownload),
		),
	}
	addAuthenticationFlags(cmd)
	addClusterDeclarationPathFlag(cmd, &opts.ClusterDeclarationPath)

	flags := cmd.Flags()
	flags.StringVarP(
		&path,
		"path",
		"p",
		"state.db",
		"determines where the downloaded state should be stored",
	)

	return cmd
}

const longMaintenanceStateDownloadDescription = `
Downloads state to the local computer. Useful for debugging the state. We recommend acquiring a lock if you intend to
modify the state
`
