package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildMaintenanceStateAcquireLockCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state-acquire-lock",
		Short: "acquires a state.db lock",
		Long:  longMaintenanceAcquireStateLockDescription,
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionMaintenanceStateAcquireLock),
			hooks.InitializeOkctl(o),
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Fprintln(o.Out, "acquiring state lock")

			err := hooks.AcquireStateLock(o)(nil, nil)
			if err != nil {
				return fmt.Errorf("acquiring state lock: %w", err)
			}

			fmt.Fprintln(o.Out, "successfully acquired lock")

			return nil
		},
		PostRunE: hooks.RunECombinator(
			hooks.EmitEndCommandExecutionEvent(metrics.ActionMaintenanceStateAcquireLock),
		),
	}

	return cmd
}

const longMaintenanceAcquireStateLockDescription = `
This command will acquire state lock. Useful when doing maintenance on a manually downloaded state.db file.
`
