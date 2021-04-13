package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/spinner"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/commands"
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
}

// Validate the options for "apply application"
func (o applyApplicationOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.File, validation.Required),
	)
}

//nolint funlen
func buildApplyApplicationCommand(o *okctl.Okctl) *cobra.Command {
	scaffoldOpts := &client.ScaffoldApplicationOpts{}
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

			scaffoldOpts.ID = &api.ID{
				Region:       o.Declaration.Metadata.Region,
				AWSAccountID: o.Declaration.Metadata.AccountID,
				ClusterName:  o.Declaration.Metadata.Name,
			}

			scaffoldOpts.Application, err = commands.InferApplicationFromStdinOrFile(o.In, o.FileSystem, opts.File)
			if err != nil {
				return fmt.Errorf("inferring application from stdin or file: %w", err)
			}

			scaffoldOpts.OutputDir, err = o.GetRepoOutputDir()
			if err != nil {
				return err
			}

			handlers := o.StateHandlers(o.StateNodes())

			hz, err := handlers.Domain.GetPrimaryHostedZone()
			if err != nil {
				return err
			}

			scaffoldOpts.HostedZoneID = hz.HostedZoneID
			scaffoldOpts.HostedZoneDomain = hz.Domain

			repo, err := handlers.Github.GetGithubRepository(
				fmt.Sprintf("%s/%s", o.Declaration.Github.Organisation, o.Declaration.Github.Repository),
			)
			if err != nil {
				return err
			}

			scaffoldOpts.IACRepoURL = repo.GitURL

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("failed validating options: %w", err)
			}

			services, _ := o.ClientServices(o.StateHandlers(o.StateNodes()))

			spin, err := spinner.New("synchronizing", o.Err)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
			}

			reconciliationManager := reconciler.NewCompositeReconciler(spin,
				reconciler.NewApplicationReconciler(services.ApplicationService),
			)

			reconciliationManager.SetCommonMetadata(&resourcetree.CommonMetadata{
				Ctx:                    o.Ctx,
				Out:                    o.Out,
				ClusterID:              *scaffoldOpts.ID,
				Declaration:            o.Declaration,
				ApplicationDeclaration: scaffoldOpts.Application,
			})

			dependencyTree := controller.CreateApplicationResourceDependencyTree()

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
