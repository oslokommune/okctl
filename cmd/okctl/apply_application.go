package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/controller/application/reconciliation"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/commands"
	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// requiredApplyApplicationArguments defines number of arguments the ApplyApplication command expects
const requiredApplyApplicationArguments = 0

// applyApplicationOpts contains all the possible options for "apply application"
type applyApplicationOpts struct {
	File string

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

			opts.Application, err = commands.InferApplicationFromStdinOrFile(*o.Declaration, o.In, o.FileSystem, opts.File)
			if err != nil {
				return fmt.Errorf(constant.InferApplicationError, err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := opts.Validate()
			if err != nil {
				return fmt.Errorf(constant.OptionValidationerror, err)
			}

			state := o.StateHandlers(o.StateNodes())

			services, _ := o.ClientServices(state)

			spin, err := spinner.New("applying application", o.Err)
			if err != nil {
				return fmt.Errorf(constant.SpinnerCreationError, err)
			}

			schedulerOpts := common.SchedulerOpts{
				Out:                             o.Out,
				Spinner:                         spin,
				ReconciliationLoopDelayFunction: common.DefaultDelayFunction,
				ClusterDeclaration:              *o.Declaration,
				ApplicationDeclaration:          opts.Application,
			}

			scheduler := common.NewScheduler(schedulerOpts,
				reconciliation.NewApplicationReconciler(services.ApplicationService),
				reconciliation.NewContainerRepositoryReconciler(services.ContainerRepository),
			)

			_, err = scheduler.Run(o.Ctx, state)
			if err != nil {
				return fmt.Errorf(constant.ReconcileApplicationError, err)
			}

			return commands.WriteApplyApplicationSuccessMessage(commands.WriteApplyApplicationSucessMessageOpts{
				Out:         o.Out,
				Application: opts.Application,
				Cluster:     *o.Declaration,
			})
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "Specify the file path. Use \"-\" for stdin")

	return cmd
}
