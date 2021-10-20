package main

import (
	"fmt"

	"github.com/oslokommune/okctl/cmd/okctl/preruns"
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
		PreRunE: preruns.PreRunECombinator(
			preruns.LoadUserData(o),
			preruns.InitializeMetrics(o),
			func(cmd *cobra.Command, args []string) error {
				metrics.Publish(metrics.Event{
					Category: metrics.CategoryCommandExecution,
					Action:   metrics.ActionVersion,
					Label:    metrics.LabelStart,
				})

				return nil
			},
		),
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := fmt.Fprint(o.Out, version.String()+"\n")
			return err
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			metrics.Publish(metrics.Event{
				Category: metrics.CategoryCommandExecution,
				Action:   metrics.ActionVersion,
				Label:    metrics.LabelEnd,
			})

			return nil
		},
	}

	return cmd
}
