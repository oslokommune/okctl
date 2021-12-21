package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildMaintenanceStateReleaseLockCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state-release-lock",
		Short: "releases a state.db lock",
		Long:  longMaintenanceReleaseStateLockDescription,
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionMaintenanceStateReleaseLock),
			hooks.InitializeOkctl(o),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Fprintln(o.Out, "releasing state lock")

			err := hooks.ReleaseStateLock(o)(nil, nil)
			if err != nil {
				return fmt.Errorf("releasing state lock: %w", err)
			}

			fmt.Fprintln(o.Out, "successfully released lock")

			return nil
		},
		PostRunE: hooks.RunECombinator(
			hooks.EmitEndCommandExecutionEvent(metrics.ActionMaintenanceStateReleaseLock),
		),
	}

	return cmd
}

const longMaintenanceReleaseStateLockDescription = `
This command will remove an existing state lock. Useful if something unexpected happened and the lock is still active.
`
