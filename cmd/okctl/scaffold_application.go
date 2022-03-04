package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
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
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionScaffoldApplication),
			func(cmd *cobra.Command, args []string) error {
				if opts.ClusterDeclarationPath != "" {
					clusterDeclaration, err := hooks.LoadClusterDeclarationFromPath(o, &opts.ClusterDeclarationPath)
					if err != nil {
						return fmt.Errorf("inferring cluster declaration: %w", err)
					}

					o.Declaration = clusterDeclaration

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
		PostRunE: hooks.RunECombinator(
			hooks.EmitEndCommandExecutionEvent(metrics.ActionScaffoldApplication),
		),
	}
	addAuthenticationFlags(cmd)
	addClusterDeclarationPathFlag(cmd, &opts.ClusterDeclarationPath)

	return cmd
}
