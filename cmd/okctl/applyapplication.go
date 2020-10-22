package main

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

// requiredApplyApplicationArguments defines number of arguments the ApplyApplication command expects
const requiredApplyApplicationArguments = 1

// applyApplicationOpts contains all the possible options for "apply application"
type applyApplicationOpts struct {
	File        string
}

// Validate the options for "apply application"
func (o *applyApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.File, validation.Required),
	)
}

func buildApplyApplicationCommand(o *okctl.Okctl) *cobra.Command {
	scaffoldOpts := &client.ScaffoldApplicationOpts{}
	opts := &applyApplicationOpts{}

	cmd := &cobra.Command{
		Use:   "application [env]",
		Short: "Applies an application.yaml to the IAC repo",
		Args:  cobra.ExactArgs(requiredApplyApplicationArguments),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			environment := args[0]

			err := o.InitialiseWithOnlyEnv(environment)
			if err != nil {
				return err
			}

			cluster := o.RepoStateWithEnv.GetCluster()
			metadata := o.RepoStateWithEnv.GetMetadata()

			scaffoldOpts.Id = &api.ID{
				Region:       o.CloudProvider.Region(),
				AWSAccountID: cluster.AWSAccountID,
				Environment:  cluster.Environment,
				Repository:   metadata.Name,
				ClusterName:  cluster.Name,
			}

			scaffoldOpts.In = o.In
			scaffoldOpts.Out = o.Out
			scaffoldOpts.ApplicationFilePath = opts.File
			scaffoldOpts.RepoDir, err = o.GetRepoDir()
			if err != nil {
				return err
			}

			hostedZone := o.RepoStateWithEnv.GetPrimaryHostedZone()
			scaffoldOpts.HostedZoneID = hostedZone.ID

			for _, repo := range cluster.Github.Repositories {
				scaffoldOpts.IACRepoURL = repo.GitURL
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("failed validating options: %w", err)
			}

			spin, _ := spinner.New("scaffolding", o.Out)
			services, _ := o.ClientServices(spin)

			err = services.ApplicationService.ScaffoldApplication(o.Ctx, scaffoldOpts)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "Specify the file path. Use \"-\" for stdin")

	return cmd
}
