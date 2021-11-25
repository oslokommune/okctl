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
				metrics.Publish(generateStartEvent(metrics.ActionScaffoldApplication))

				if declarationPath != "" {
					clusterDeclaration, err := commands.InferClusterFromStdinOrFile(o.In, declarationPath)
					if err != nil {
						return fmt.Errorf("inferring cluster declaration: %w", err)
					}

					o.Declaration = clusterDeclaration

					err = commands.ValidateBinaryVersionNotLessThanClusterVersion(o)
					if err != nil {
						return err
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
			metrics.Publish(generateEndEvent(metrics.ActionScaffoldApplication))

			return nil
		},
	}

	return cmd
}
