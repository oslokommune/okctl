package main

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/controller"
	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// requiredApplyApplicationArguments defines number of arguments the ApplyApplication command expects
const requiredApplyApplicationArguments = 1

// applyApplicationOpts contains all the possible options for "apply application"
type applyApplicationOpts struct {
	File string
}

// Validate the options for "apply application"
func (o *applyApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.File, validation.Required),
	)
}

//nolint funlen
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
				errEnvNotFound := &okctl.ErrorEnvironmentNotFound{}

				if errors.As(err, &errEnvNotFound) {
					fmt.Fprintf(o.Err, "\nThe specified environment \"%s\" does not exist.", errEnvNotFound.TargetEnvironment)
					fmt.Fprintf(o.Err, "Available environments: %v\n", errEnvNotFound.AvailableEnvironments)
				}

				return err
			}

			cluster := o.RepoStateWithEnv.GetCluster()
			metadata := o.RepoStateWithEnv.GetMetadata()

			scaffoldOpts.ID = &api.ID{
				Region:       o.CloudProvider.Region(),
				AWSAccountID: cluster.AWSAccountID,
				Environment:  cluster.Environment,
				Repository:   metadata.Name,
				ClusterName:  o.RepoStateWithEnv.GetClusterName(),
			}

			scaffoldOpts.Application, err = commands.InferApplicationFromStdinOrFile(o.In, o.FileSystem, opts.File)
			if err != nil {
				return fmt.Errorf("inferring application from stdin or file: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("failed validating options: %w", err)
			}

			spin, _ := spinner.New("synchronizing ", o.Out)
			services, _ := o.ClientServices(spin)

			reconciliationManager := reconciler.NewCompositeReconciler(spin,
				reconciler.NewApplicationReconciler(services.ApplicationService),
			)

			reconciliationManager.SetCommonMetadata(&resourcetree.CommonMetadata{
				Ctx:       o.Ctx,
				Out:       o.Out,
				ClusterID: *scaffoldOpts.ID,
				// TODO: should pass in cluster as done in the huge rewrite PR
				Declaration: &v1alpha1.Cluster{
					Github: v1alpha1.ClusterGithub{
						Repository: commands.GetFirstGithubRepositoryURL(o.RepoStateWithEnv.GetGithub().Repositories),
						OutputPath: o.RepoStateWithEnv.GetMetadata().OutputDir,
					},
				},
			})

			dependencyTree := controller.CreateApplicationResourceDependencyTree()

			dependencyTree.SetStateRefresher(resourcetree.ResourceNodeTypeApplication, func(node *resourcetree.ResourceNode) {
				primaryHostedZone := o.RepoStateWithEnv.GetPrimaryHostedZone()

				node.ResourceState = &reconciler.ApplicationState{
					Declaration: scaffoldOpts.Application,

					PrimaryHostedZoneID:     primaryHostedZone.ID,
					PrimaryHostedZoneDomain: primaryHostedZone.Domain,
				}
			})

			err = commands.SynchronizeApplication(reconciliationManager, dependencyTree)
			if err != nil {
				return fmt.Errorf("synchronizing application: %w", err)
			}

			return commands.WriteApplyApplicationSuccessMessage(
				o.Out,
				scaffoldOpts.Application.Name,
				scaffoldOpts.OutputDir,
			)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "Specify the file path. Use \"-\" for stdin")

	return cmd
}
