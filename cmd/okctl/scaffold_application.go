package main

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/scaffold"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const requiredArgumentsForCreateApplicationCommand = 1

// scaffoldApplicationOpts contains all the possible options for "scaffold application"
type scaffoldApplicationOpts struct {
	Environment string
	Outfile     string
}

// Validate the options for "create application"
func (o *scaffoldApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment),
	)
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
			opts.Environment = args[0]

			err := o.InitialiseWithOnlyEnv(opts.Environment)
			if err != nil {
				return err
			}

			hostedZone := o.RepoStateWithEnv.GetPrimaryHostedZone()
			interpolationOpts.PrimaryHostedZone = hostedZone.Domain

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Validate()
			if err != nil {
				return err
			}

			okctlAppTemplate, err := scaffold.GenerateOkctlAppTemplate(interpolationOpts)
			if err != nil {
				return err
			}

			err = scaffold.SaveOkctlAppTemplate(o.FileSystem, opts.Outfile, okctlAppTemplate)
			if err != nil {
				return err
			}

			fmt.Fprintln(o.Out, "Scaffolding successful.")
			fmt.Fprintf(
				o.Out,
				"Edit %s to your liking and run %s\n",
				opts.Outfile,
				aurora.Green(fmt.Sprintf("okctl apply application %s -f %s\n", opts.Environment, opts.Outfile)),
			)

			return err
		},
	}

	cmd.Flags().StringVarP(&opts.Outfile, "out", "o", "application.yaml", "specify where to output result")

	return cmd
}
