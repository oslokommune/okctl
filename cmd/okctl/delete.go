package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/request"
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
			opts.RepositoryName = o.RepoData.Name

			err := opts.Validate()
			if err != nil {
				return err
			}

			awsAccountID, err := o.AWSAccountID(opts.Environment)
			if err != nil {
				return err
			}

			return o.Initialise(opts.Environment, awsAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			data, err := json.Marshal(opts)
			if err != nil {
				return err
			}

			r := request.New(fmt.Sprintf("http://%s/v1/", o.Destination))

			resp, err := r.Delete("clusters/", data)
			if err != nil {
				return err
			}

			_, err = io.Copy(o.Out, strings.NewReader(resp))
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
