package main

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

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
		Short: "Scaffold an application template",
		Long:  "Scaffolds an application.yaml template which can be used to produce necessary Kubernetes and ArgoCD resources",
		Args:  cobra.ExactArgs(requiredArgumentsForCreateApplicationCommand),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if declarationPath != "" {
				clusterDeclaration, err := commands.InferClusterFromStdinOrFile(o.In, declarationPath)
				if err != nil {
					return fmt.Errorf(constant.InferClusterDeclarationError, err)
				}

				opts.PrimaryHostedZone = clusterDeclaration.ClusterRootDomain
			} else {
				opts.PrimaryHostedZone = "okctl.io"
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ScaffoldApplicationDeclaration(o.Out, opts)
		},
	}

	return cmd
}
