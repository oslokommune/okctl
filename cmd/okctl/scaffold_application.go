package main

import (
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := o.Initialise()
			if err != nil {
				return err
			}

			hostedZone, err := o.StateHandlers(o.StateNodes()).Domain.GetPrimaryHostedZone()
			if err != nil {
				return err
			}

			opts.PrimaryHostedZone = hostedZone.Domain

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ScaffoldApplicationDeclaration(o.Out, opts)
		},
	}

	return cmd
}
