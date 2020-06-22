package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/oslokommune/okctl/pkg/credentials/login"
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
	var opts okctl.DeleteClusterOpts

	cmd := &cobra.Command{
		Use:   "cluster [env]",
		Short: "Delete a cluster",
		Long: `Delete all resources related to an EKS cluster,
including VPC, this is a highly destructive operation.`,
		Args: cobra.ExactArgs(deleteClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Environment = args[0]

			err := opts.Valid()
			if err != nil {
				return err
			}

			if o.NoInput {
				return fmt.Errorf("delete cluster requires user input for now")
			}

			awsAccountID, err := o.AWSAccountID(opts.Environment)
			if err != nil {
				return err
			}

			l, err := login.Interactive(awsAccountID, o.Region(), o.Username())
			if err != nil {
				return err
			}

			o.CredentialsProvider = credentials.New(l)

			c, err := cloud.New(o.Region(), o.CredentialsProvider)
			if err != nil {
				return err
			}

			o.CloudProvider = c.Provider

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return o.DeleteCluster(opts)
		},
		PostRunE: func(_ *cobra.Command, args []string) error {
			return o.WriteCurrentRepoData()
		},
	}

	return cmd
}
