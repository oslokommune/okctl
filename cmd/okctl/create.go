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
	createClusterArgs = 2
	defaultCidr       = "192.168.0.0/20"
)

func buildCreateCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create commands",
	}

	cmd.AddCommand(buildCreateClusterCommand(o))

	return cmd
}

// nolint: funlen
func buildCreateClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &api.ClusterCreateOpts{}

	cmd := &cobra.Command{
		Use:   "cluster [env] [AWS account id]",
		Short: "Create a cluster",
		Long: `Fetch all tasks required to get an EKS cluster up and running on AWS.
This includes creating an EKS compatible VPC with private, public
and database subnets.`,
		Args: cobra.ExactArgs(createClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Environment = args[0]
			opts.AWSAccountID = args[1]
			opts.RepositoryName = o.RepoData.Name
			opts.ClusterName = o.ClusterName(opts.Environment)
			opts.Region = o.Region()

			err := opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate create cluster options", errors.Invalid)
			}

			return o.Initialise(opts.Environment, opts.AWSAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			// Discarding the output for now, until we restructure
			// the API to return everything we need to write
			// the result ourselves
			c := client.New(ioutil.Discard, o.ServerURL)

			err := c.CreateVpc(&api.CreateVpcOpts{
				AwsAccountID: opts.AWSAccountID,
				ClusterName:  opts.ClusterName,
				Env:          opts.Environment,
				RepoName:     opts.RepositoryName,
				Cidr:         opts.Cidr,
				Region:       opts.Region,
			})
			if err != nil {
				return err
			}

			err = c.CreateClusterConfig(&api.CreateClusterConfigOpts{
				ClusterName:  opts.ClusterName,
				Region:       opts.Region,
				Cidr:         opts.Cidr,
				AwsAccountID: opts.AWSAccountID,
			})
			if err != nil {
				return err
			}

			err = c.CreateCluster(opts)
			if err != nil {
				return err
			}

			policy, err := c.CreateExternalSecretsPolicy(&api.CreateExternalSecretsPolicyOpts{
				Repository:  opts.RepositoryName,
				Environment: opts.Environment,
			})
			if err != nil {
				return err
			}

			return c.CreateExternalSecretsServiceAccount(&api.CreateExternalSecretsServiceAccountOpts{
				ClusterName:  opts.ClusterName,
				Environment:  opts.Environment,
				Region:       opts.Region,
				AWSAccountID: opts.AWSAccountID,
				PolicyArn:    policy.PolicyARN,
			})
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Cidr, "cidr", "c", defaultCidr,
		"CIDR block the AWS VPC and subnets are created within")

	return cmd
}
