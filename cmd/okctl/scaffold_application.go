package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/cmd/okctl/preruns"
	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const requiredArgumentsForCreateApplicationCommand = 0

// nolint: funlen
func buildScaffoldApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := commands.ScaffoldApplicationOpts{}

	cmd := &cobra.Command{
		Use:   "application",
		Short: ScaffoldShortDescription,
		Long:  ScaffoldLongDescription,
		Args:  cobra.ExactArgs(requiredArgumentsForCreateApplicationCommand),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PreRunE: preruns.PreRunECombinator(
			preruns.LoadUserData(o),
			preruns.InitializeMetrics(o),
			func(cmd *cobra.Command, args []string) error {
				metrics.Publish(metrics.Event{
					Category: metrics.CategoryCommandExecution,
					Action:   metrics.ActionScaffoldApplication,
					Label:    metrics.LabelStart,
				})

				if declarationPath != "" {
					clusterDeclaration, err := commands.InferClusterFromStdinOrFile(o.In, declarationPath)
					if err != nil {
						return fmt.Errorf("inferring cluster declaration: %w", err)
					}

					opts.PrimaryHostedZone = clusterDeclaration.ClusterRootDomain
				} else {
					opts.PrimaryHostedZone = "okctl.io"
				}

				return nil
			},
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ScaffoldApplicationDeclaration(o.Out, opts)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			metrics.Publish(metrics.Event{
				Category: metrics.CategoryCommandExecution,
				Action:   metrics.ActionScaffoldApplication,
				Label:    metrics.LabelEnd,
			})

			return nil
		},
	}

	return cmd
}
