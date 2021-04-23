package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/spinner"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/controller"
	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// requiredApplyApplicationArguments defines number of arguments the ApplyApplication command expects
const requiredApplyApplicationArguments = 0

// applyApplicationOpts contains all the possible options for "apply application"
type applyApplicationOpts struct {
	File string

	ClusterID   api.ID
	Application v1alpha1.Application
}

// Validate the options for "apply application"
func (o applyApplicationOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.File, validation.Required),
		validation.Field(&o.Application),
	)
}

//nolint funlen
func buildApplyApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := &applyApplicationOpts{}

	cmd := &cobra.Command{
		Use:   "application",
		Short: "Applies an application.yaml to the IAC repo",
		Args:  cobra.ExactArgs(requiredApplyApplicationArguments),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := o.Initialise()
			if err != nil {
				return err
			}

			opts.ClusterID = api.ID{
				Region:       o.Declaration.Metadata.Region,
				AWSAccountID: o.Declaration.Metadata.AccountID,
				ClusterName:  o.Declaration.Metadata.Name,
			}

			opts.Application, err = commands.InferApplicationFromStdinOrFile(*o.Declaration, o.In, o.FileSystem, opts.File)
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

			handlers := o.StateHandlers(o.StateNodes())

			services, _ := o.ClientServices(handlers)

			spin, err := spinner.New("synchronizing", o.Err)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
			}

			reconciliationManager := reconciler.NewCompositeReconciler(spin,
				reconciler.NewApplicationReconciler(services.ApplicationService),
				reconciler.NewContainerRepositoryReconciler(services.ContainerRepository),
			)

			reconciliationManager.SetCommonMetadata(&resourcetree.CommonMetadata{
				Ctx:                    o.Ctx,
				Out:                    o.Out,
				ClusterID:              opts.ClusterID,
				Declaration:            o.Declaration,
				ApplicationDeclaration: opts.Application,
			})

			reconciliationManager.SetStateHandlers(handlers)

			dependencyTree := controller.CreateApplicationResourceDependencyTree()

			err = commands.SynchronizeApplication(commands.SynchronizeApplicationOpts{
				ReconciliationManager: reconciliationManager,
				Application:           opts.Application,
				Tree:                  dependencyTree,
			})
			if err != nil {
				return fmt.Errorf("synchronizing application: %w", err)
			}

			return commands.WriteApplyApplicationSuccessMessage(
				o.Out,
				opts.Application.Metadata.Name,
				o.Declaration.Github.OutputPath,
			)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "Specify the file path. Use \"-\" for stdin")

	return cmd
}
