package main

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/scaffold"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const requiredArgumentsForCreateApplicationCommand = 0

// scaffoldApplicationOpts contains all the possible options for "scaffold application"
type scaffoldApplicationOpts struct {
	Outfile string
}

// nolint: funlen
func buildCreateApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := &scaffoldApplicationOpts{}
	interpolationOpts := &scaffold.InterpolationOpts{}

	cmd := &cobra.Command{
		Use:   "application ENV",
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

			interpolationOpts.PrimaryHostedZone = hostedZone.Domain

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			okctlAppTemplate, err := scaffold.GenerateOkctlAppTemplate(interpolationOpts)
			if err != nil {
				return err
			}

			err = scaffold.SaveOkctlAppTemplate(o.FileSystem, opts.Outfile, okctlAppTemplate)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintln(o.Out, "Scaffolding successful.")
			_, _ = fmt.Fprintf(
				o.Out,
				"Edit %s to your liking and run %s\n",
				opts.Outfile,
				aurora.Green(fmt.Sprintf("okctl apply application -f %s\n", opts.Outfile)),
			)

			return err
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "out", "o", "application.yaml", "specify where to output result")

	return cmd
}
