package main

import (
	"fmt"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func buildVersionCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: VersionShortDescription,
		Long:  VersionLongDescription,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionVersion),
		),
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := fmt.Fprint(o.Out, version.String()+"\n")
			return err
		},
		PostRunE: hooks.RunECombinator(
			hooks.EmitEndCommandExecutionEvent(metrics.ActionVersion),
		),
	}

	return cmd
}
