package main

import (
	"fmt"

	"github.com/logrusorgru/aurora"

	kaex "github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/scaffold"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const requiredArgumentsForCreateApplicationCommand = 1

// CreateApplicationOpts contains all the possible options for "create application"
type CreateApplicationOpts struct {
	Environment string
	KaexOpts    *kaex.Kaex
	Outfile     string
}

// Validate the options for "create application"
func (o *CreateApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment),
	)
}

// nolint: funlen
func buildCreateApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := &CreateApplicationOpts{}
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
			interpolationOpts.Domain = hostedZone.Domain

			opts.KaexOpts = &kaex.Kaex{
				Err:             o.Err,
				Out:             o.Out,
				In:              o.In,
				TemplatesDirURL: "https://raw.githubusercontent.com/oslokommune/kaex/master/templates",
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Validate()
			if err != nil {
				return err
			}

			template, err := scaffold.FetchTemplate(*opts.KaexOpts)
			if err != nil {
				return fmt.Errorf("failed fetching application.yaml example: %w", err)
			}

			interpolatedResult, err := scaffold.InterpolateTemplate(template, interpolationOpts)
			if err != nil {
				return err
			}

			err = scaffold.SaveTemplate(opts.Outfile, interpolatedResult)
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
