package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/request"
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
			r := request.New(fmt.Sprintf("http://%s/v1/", o.Destination))

			{
				vpcOpts := &api.CreateVpcOpts{
					AwsAccountID: opts.AWSAccountID,
					ClusterName:  opts.ClusterName,
					Env:          opts.Environment,
					RepoName:     opts.RepositoryName,
					Cidr:         opts.Cidr,
					Region:       opts.Region,
				}

				vpcData, err := json.Marshal(vpcOpts)
				if err != nil {
					return errors.E(err, "failed to marshal create vpc request")
				}

				resp, err := r.Post("vpcs/", vpcData)
				if err != nil {
					return errors.E(err, resp, errors.Internal)
				}

				_, err = io.Copy(o.Out, strings.NewReader(resp))
				if err != nil {
					return err
				}
			}

			{
				cfgOpts := &api.CreateClusterConfigOpts{
					ClusterName:  opts.ClusterName,
					Region:       opts.Region,
					Cidr:         opts.Cidr,
					AwsAccountID: opts.AWSAccountID,
				}

				cfgData, err := json.Marshal(cfgOpts)
				if err != nil {
					return errors.E(err, "failed to marshal create cluster config request")
				}

				resp, err := r.Post("vpcs/", cfgData)
				if err != nil {
					return errors.E(err, resp, errors.Internal)
				}

				_, err = io.Copy(o.Out, strings.NewReader(resp))
				if err != nil {
					return err
				}
			}

			{
				data, err := json.Marshal(opts)
				if err != nil {
					return errors.E(err, "failed to marshal create cluster request")
				}

				resp, err := r.Post("clusters/", data)
				if err != nil {
					return errors.E(err, resp, errors.Internal)
				}

				_, err = io.Copy(o.Out, strings.NewReader(resp))
				if err != nil {
					return err
				}

				return nil
			}
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Cidr, "cidr", "c", defaultCidr,
		"CIDR block the AWS VPC and subnets are created within")

	return cmd
}
