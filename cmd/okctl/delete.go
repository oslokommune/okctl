package main

import (
	"io/ioutil"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
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

// nolint: funlen
func buildDeleteClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &api.ClusterDeleteOpts{}

	cmd := &cobra.Command{
		Use:   "cluster [env]",
		Short: "Delete a cluster",
		Long: `Delete all resources related to an EKS cluster,
including VPC, this is a highly destructive operation.`,
		Args: cobra.ExactArgs(deleteClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Environment = args[0]
			opts.RepositoryName = o.RepoData.Name

			err := opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate delete cluster options")
			}

			awsAccountID, err := o.AWSAccountID(opts.Environment)
			if err != nil {
				return err
			}

			return o.Initialise(opts.Environment, awsAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			// Discarding the output for now until we have
			// restructured the API and handle the response
			// properly
			c := client.New(o.Debug, ioutil.Discard, o.ServerURL)

			err := c.DeleteCluster(opts)
			if err != nil {
				return err
			}

			return c.DeleteVpc(&api.DeleteVpcOpts{
				Env:      opts.Environment,
				RepoName: opts.RepositoryName,
			})
		},
	}

	return cmd
}
