package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/okctlapplication"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// CreateApplicationOpts contains all the possible options for "create application"
type CreateApplicationOpts struct {
	fullExample bool
	env         string
}

// Validate the options for "create application"
func (o *CreateApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.fullExample),
		validation.Field(&o.env),
	)
}

// nolint funlen
func buildCreateApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := &CreateApplicationOpts{}

	cmd := &cobra.Command{
		Use:   "application",
		Short: "Create an application template",
		Long:  "Scaffolds an application.yaml template which can be used to produce necessary Kubernetes and ArgoCD resources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var buffer bytes.Buffer

			if opts.fullExample {
				err := api.FetchFullExample(&buffer)
				if err != nil {
					return fmt.Errorf("failed fetching application.yaml example: %w", err)
				}
			} else {
				err := api.FetchMinimalExample(&buffer)
				if err != nil {
					return fmt.Errorf("failed fetching application.yaml example: %w", err)
				}
			}

			cluster := okctlapplication.GetCluster(o, cmd, opts.env)
			if cluster == nil {
				fmt.Fprint(o.Out, buffer)

				return nil
			}

			output := strings.Replace(
				buffer.String(),
				"my-domain.io",
				fmt.Sprintf("<app-name>.%s-%s.oslo.systems", cluster.Name, cluster.Environment),
				1,
			)

			fmt.Fprint(o.Out, output)

			return nil
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.fullExample, "full", "f", false, "Scaffold a full rather than a minimal example")
	flags.StringVarP(&opts.env, "environment", "e", "", "Use a certain environment as base for the scaffold")

	return cmd
}
