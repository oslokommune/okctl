package main

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4/is"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

const (
	addUserArgs = 2
)

// AddUserOpts options for the adduser command
type AddUserOpts struct {
	Environment    string
	AWSAccountID   string
	RepositoryName string
	Region         string
	ClusterName    string
	UserPoolID     string
	UserEmail      string
}

// Validate the inputs
func (o *AddUserOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("must consist of 3-64 characters (a-z, A-Z)")),
		validation.Field(&o.AWSAccountID, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{12}$")).Error("must consist of 12 digits")),
		validation.Field(&o.RepositoryName, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.UserPoolID, validation.Required),
		validation.Field(&o.UserEmail, validation.Required, is.Email),
	)
}

// nolint: funlen
func buildAddUserCommand(o *okctl.Okctl) *cobra.Command {
	opts := &AddUserOpts{}
	cmd := &cobra.Command{
		Use:   "adduser [env] [email]",
		Short: "Add a user to identitypool",
		Long:  `Add a user to the identitypool associated with the cluster, for login in ArgoCD and other applications`,
		Args:  cobra.ExactArgs(addUserArgs),
		RunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]
			userEmail := args[1]

			err := validation.Validate(
				&environment,
				validation.Required,
				validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("the environment must consist of 3-64 characters (a-z, A-Z)"),
			)
			if err != nil {
				return err
			}

			err = o.InitialiseWithOnlyEnv(environment)
			if err != nil {
				return err
			}

			meta := o.RepoStateWithEnv.GetMetadata()
			cluster := o.RepoStateWithEnv.GetCluster()

			opts.Environment = environment
			opts.AWSAccountID = cluster.AWSAccountID
			opts.RepositoryName = meta.Name
			opts.Region = meta.Region
			opts.UserPoolID = cluster.IdentityPool.UserPoolID
			opts.ClusterName = cluster.Name
			opts.UserEmail = userEmail

			err = opts.Validate()
			if err != nil {
				return err
			}

			id := api.ID{
				Region:       meta.Region,
				AWSAccountID: opts.AWSAccountID,
				Environment:  opts.Environment,
				Repository:   opts.RepositoryName,
				ClusterName:  opts.ClusterName,
			}

			spin, err := spinner.New("creating-identity-pool-users", o.Err)
			if err != nil {
				return err
			}

			services, err := o.ClientServices(spin)
			if err != nil {
				return nil
			}

			_, err = services.IdentityManager.CreateIdentityPoolUser(o.Ctx, api.CreateIdentityPoolUserOpts{
				ID:         id,
				Email:      opts.UserEmail,
				UserPoolID: opts.UserPoolID,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
