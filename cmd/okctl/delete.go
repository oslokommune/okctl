package main

import (
	"io/ioutil"
	"path"

	"github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/client/core/api/rest"
	"github.com/oslokommune/okctl/pkg/client/core/store/filesystem"
	"github.com/oslokommune/okctl/pkg/config"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	deleteClusterArgs = 1
)

func buildDeleteCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete commands",
	}

	cmd.AddCommand(buildDeleteClusterCommand(o))

	return cmd
}

// DeleteClusterOpts contains the required inputs
type DeleteClusterOpts struct {
	Region       string
	AWSAccountID string
	Environment  string
	Repository   string
	ClusterName  string
}

// Validate the inputs
func (o *DeleteClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: funlen
func buildDeleteClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &DeleteClusterOpts{}

	cmd := &cobra.Command{
		Use:   "cluster [env]",
		Short: "Delete a cluster",
		Long: `Delete all resources related to an EKS cluster,
including VPC, this is a highly destructive operation.`,
		Args: cobra.ExactArgs(deleteClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]

			opts.Region = o.Region()
			opts.AWSAccountID = o.AWSAccountID(environment)
			opts.Environment = environment
			opts.Repository = o.RepoData.Name
			opts.ClusterName = o.ClusterName(environment)

			err := opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate delete cluster options")
			}

			return o.Initialise(opts.Environment, opts.AWSAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			// Discarding the output for now until we have
			// restructured the API and handle the response
			// properly
			c := rest.New(o.Debug, ioutil.Discard, o.ServerURL)

			id := api.ID{
				Region:       opts.Region,
				AWSAccountID: opts.AWSAccountID,
				Environment:  opts.Environment,
				Repository:   opts.Repository,
				ClusterName:  opts.ClusterName,
			}

			outputDir, err := o.GetRepoOutputDir(opts.Environment)
			if err != nil {
				return err
			}

			repoDir, err := o.GetRepoDir()
			if err != nil {
				return err
			}

			clusterService := core.NewClusterService(
				rest.NewClusterAPI(c),
				filesystem.NewClusterStore(
					filesystem.Paths{
						ConfigFile: config.DefaultRepositoryConfig,
						BaseDir:    repoDir,
					},
					filesystem.Paths{
						ConfigFile: config.DefaultClusterConfig,
						BaseDir:    path.Join(outputDir, config.DefaultClusterBaseDir),
					},
					o.FileSystem,
					o.RepoData,
				),
			)

			err = clusterService.DeleteCluster(o.Ctx, api.ClusterDeleteOpts{
				ID: id,
			})
			if err != nil {
				return err
			}

			vpcService := core.NewVPCService(
				rest.NewVPCAPI(c),
				filesystem.NewVpcStore(
					config.DefaultVpcOutputs,
					config.DefaultVpcCloudFormationTemplate,
					path.Join(outputDir, config.DefaultVpcBaseDir),
					o.FileSystem,
				),
			)

			err = vpcService.DeleteVpc(o.Ctx, api.DeleteVpcOpts{
				ID: id,
			})

			return nil
		},
	}

	return cmd
}
