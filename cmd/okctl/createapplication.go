package main

import (
	"fmt"

	kaex "github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/scaffold"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// CreateApplicationOpts contains all the possible options for "create application"
type CreateApplicationOpts struct {
	Environment string
}

// Validate the options for "create application"
func (o *CreateApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment),
	)
}

func buildCreateApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := &CreateApplicationOpts{}

	cmd := &cobra.Command{
		Use:   "application",
		Short: "Create an application template",
		Long:  "Scaffolds an application.yaml template which can be used to produce necessary Kubernetes and ArgoCD resources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			kx := kaex.Kaex{
				Err:             o.Err,
				Out:             o.Out,
				In:              o.In,
				TemplatesDirURL: "https://raw.githubusercontent.com/oslokommune/kaex/master/templates",
			}

			template, err := scaffold.FetchTemplate(kx)
			if err != nil {
				return fmt.Errorf("failed fetching application.yaml example: %w", err)
			}

			interpolatedResult, err := scaffold.InterpolateTemplate(o, cmd, opts.Environment, template)
			if err != nil {
				return err
			}

			err = scaffold.SaveTemplate(interpolatedResult)
			if err != nil {
				return err
			}

			fmt.Fprintln(o.Out, "Scaffolding successful.")
			fmt.Fprintln(o.Out, "Edit ./application.yaml to your liking and run okctl apply application -f application.yaml")

			return err
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Environment, "environment", "e", "", "Use a certain environment as base for the scaffold")

	return cmd
}
