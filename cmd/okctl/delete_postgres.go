package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type deletePostgresOpts struct {
	ID              api.ID
	Environment     string
	ApplicationName string
	VpcID           string
}

// Validate the inputs
func (o *deletePostgresOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.ApplicationName, validation.Required),
		validation.Field(&o.VpcID, validation.Required),
	)
}

// nolint: funlen
func buildDeletePostgresCommand(o *okctl.Okctl) *cobra.Command {
	opts := &deletePostgresOpts{}

	cmd := &cobra.Command{
		Use:   "postgres",
		Short: "Delete an AWS RDS Postgres database",
		Long:  `Delete the AWS RDS Postgres database`,
		Args:  cobra.ExactArgs(0), // nolint: gomnd
		PreRunE: func(_ *cobra.Command, _ []string) error {
			err := o.InitialiseWithOnlyEnv(opts.Environment)
			if err != nil {
				return err
			}

			meta := o.RepoStateWithEnv.GetMetadata()
			cluster := o.RepoStateWithEnv.GetCluster()

			opts.ID.Environment = opts.Environment
			opts.ID.AWSAccountID = cluster.AWSAccountID
			opts.ID.Repository = meta.Name
			opts.ID.Region = meta.Region
			opts.ID.ClusterName = o.RepoStateWithEnv.GetClusterName()
			opts.VpcID = o.RepoStateWithEnv.GetVPC().VpcID

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

	flags.StringVarP(&opts.Environment,
		"environment",
		"e",
		"",
		"The environment the postgres database was created in",
	)

	flags.StringVarP(&opts.ApplicationName,
		"name",
		"n",
		"",
		"The name of the database to delete",
	)

	return cmd
}
