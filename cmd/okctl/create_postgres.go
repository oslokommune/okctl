package main

import (
	"log"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/sanity-io/litter"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	defaultPostgresUserName = "administrator"
)

type createPostgresOpts struct {
	ID              api.ID
	ApplicationName string
	Namespace       string
	UserName        string
}

// Validate the inputs
func (o *createPostgresOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ApplicationName, validation.Required),
		validation.Field(&o.UserName, validation.Required, validation.NotIn("admin")),
		validation.Field(&o.Namespace, validation.Required),
	)
}

// nolint: funlen
func buildCreatePostgresCommand(o *okctl.Okctl) *cobra.Command {
	opts := &createPostgresOpts{}

	cmd := &cobra.Command{
		Use:   "postgres APPLICATION_NAME NAMESPACE",
		Short: "Create an AWS RDS Postgres database",
		Long: `We will create an AWS RDS Postgres database and make a Secret and ConfigMap available
in the provided namespace containing data for accessing the database.
`,
		Args: cobra.ExactArgs(2), // nolint: gomnd
		PreRunE: func(cmd *cobra.Command, args []string) error {
			applicationName := args[0]
			namespace := args[1]

			err := o.Initialise()
			if err != nil {
				return err
			}

			opts.ID.AWSAccountID = o.Declaration.Metadata.AccountID
			opts.ID.Region = o.Declaration.Metadata.Region
			opts.ID.ClusterName = o.Declaration.Metadata.Name
			opts.ApplicationName = applicationName
			opts.Namespace = namespace

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

			vpc, err := o.StateHandlers(o.StateNodes()).Vpc.GetVpc(
				cfn.NewStackNamer().Vpc(o.Declaration.Metadata.Name),
			)
			if err != nil {
				return err
			}

			var ids, cidrs []string

			for _, s := range vpc.DatabaseSubnets {
				ids = append(ids, s.ID)
				cidrs = append(cidrs, s.Cidr)
			}

			db, err := services.Component.CreatePostgresDatabase(o.Ctx, client.CreatePostgresDatabaseOpts{
				ID:                opts.ID,
				ApplicationName:   opts.ApplicationName,
				UserName:          opts.UserName,
				VpcID:             vpc.VpcID,
				DBSubnetGroupName: vpc.DatabaseSubnetsGroupName,
				DBSubnetIDs:       ids,
				DBSubnetCIDRs:     cidrs,
				Namespace:         opts.Namespace,
			})
			if err != nil {
				return err
			}

			log.Println(litter.Sdump(db))

			return nil
		},
		Hidden: true,
	}

	f := cmd.Flags()
	f.StringVarP(&opts.UserName, "username", "u", defaultPostgresUserName,
		"Username to give to the postgres database admin account")

	return cmd
}
