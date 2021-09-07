package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type deletePostgresOpts struct {
	ID              api.ID
	ApplicationName string
	Namespace       string
	VpcID           string
}

// Validate the inputs
func (o *deletePostgresOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ApplicationName, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
		validation.Field(&o.VpcID, validation.Required),
	)
}

// nolint: funlen
func buildDeletePostgresCommand(o *okctl.Okctl) *cobra.Command {
	opts := &deletePostgresOpts{}

	cmd := &cobra.Command{
		Use:   "postgres",
		Short: DeletePostgresShortDescription,
		Long:  DeletePostgresLongDescription,
		Args:  cobra.ExactArgs(0), // nolint: gomnd
		PreRunE: func(_ *cobra.Command, _ []string) error {
			err := o.Initialise()
			if err != nil {
				return err
			}

			opts.ID.AWSAccountID = o.Declaration.Metadata.AccountID
			opts.ID.Region = o.Declaration.Metadata.Region
			opts.ID.ClusterName = o.Declaration.Metadata.Name

			for _, db := range o.Declaration.Databases.Postgres {
				if db.Name == opts.ApplicationName {
					opts.Namespace = db.Namespace
				}
			}

			vpc, err := o.StateHandlers(o.StateNodes()).Vpc.GetVpc(
				cfn.NewStackNamer().Vpc(o.Declaration.Metadata.Name),
			)
			if err != nil {
				return err
			}

			opts.VpcID = vpc.VpcID

			err = opts.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
			if err != nil {
				return err
			}

			err = services.Component.DeletePostgresDatabase(o.Ctx, client.DeletePostgresDatabaseOpts{
				ID:              opts.ID,
				Namespace:       opts.Namespace,
				ApplicationName: opts.ApplicationName,
				VpcID:           opts.VpcID,
			})
			if err != nil {
				return err
			}

			return nil
		},
		Hidden: true,
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.ApplicationName,
		"name",
		"n",
		"",
		"The name of the database to delete",
	)

	return cmd
}
