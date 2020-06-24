package main

import (
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

			err := opts.Validate()
			if err != nil {
				return err
			}

			awsAccountID, err := o.AWSAccountID(opts.Environment)
			if err != nil {
				return err
			}

			return o.NewProviders(opts.Environment, awsAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return o.DeleteCluster(opts)
		},
	}

	return cmd
}
