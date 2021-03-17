package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

type deletePostgresOpts struct {
	ID              api.ID
	ApplicationName string
	VpcID           string
}

// Validate the inputs
func (o *deletePostgresOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ApplicationName, validation.Required),
		validation.Field(&o.VpcID, validation.Required),
	)
}

// nolint: funlen
func buildDeletePostgresCommand(o *okctl.Okctl) *cobra.Command {
	opts := &deletePostgresOpts{}

	cmd := &cobra.Command{
		Use:   "postgres ENV APPLICATION_NAME",
		Short: "Delete an AWS RDS Postgres database",
		Long:  `Delete the AWS RDS Postgres database`,
		Args:  cobra.ExactArgs(2), // nolint: gomnd
		PreRunE: func(cmd *cobra.Command, args []string) error {
			environment := args[0]
			applicationName := args[1]

			err := o.InitialiseWithOnlyEnv(environment)
			if err != nil {
				return err
			}

			meta := o.RepoStateWithEnv.GetMetadata()
			cluster := o.RepoStateWithEnv.GetCluster()

			opts.ID.Environment = environment
			opts.ID.AWSAccountID = cluster.AWSAccountID
			opts.ID.Repository = meta.Name
			opts.ID.Region = meta.Region
			opts.ID.ClusterName = cluster.Name
			opts.ApplicationName = applicationName
			opts.VpcID = o.RepoStateWithEnv.GetVPC().VpcID

			err = opts.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spin, err := spinner.New("deleting", o.Err)
			if err != nil {
				return err
			}

			services, err := o.ClientServices(spin)
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

	return cmd
}
