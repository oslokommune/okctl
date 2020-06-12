package main

import (
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

func buildCreateClusterCommand(o *okctl.Okctl) *cobra.Command {
	var opts okctl.CreateClusterOpts

	cmd := &cobra.Command{
		Use:   "cluster [env] [AWS account id]",
		Short: "Create a cluster",
		Long: `Run all tasks required to get an EKS cluster up and running on AWS.
This includes creating an EKS compatible VPC with private, public
and database subnets.`,
		Args: cobra.ExactArgs(createClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Environment = args[0]
			opts.AWSAccountID = args[1]
			return opts.Valid()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return o.CreateCluster(opts)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Cidr, "cidr", "c", defaultCidr,
		"CIDR block the AWS VPC and subnets are created within")

	return cmd
}
