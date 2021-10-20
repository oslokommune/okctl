package main

import (
	"fmt"

	"github.com/oslokommune/okctl/cmd/okctl/preruns"
	"github.com/oslokommune/okctl/pkg/metrics"

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
		Short: ApplyApplicationShortDescription,
		Args:  cobra.ExactArgs(requiredApplyApplicationArguments),
		PreRunE: preruns.PreRunECombinator(
			preruns.InitializeMetrics(o),
			preruns.InitializeOkctl(o),
			func(cmd *cobra.Command, args []string) (err error) {
				metrics.Publish(metrics.Event{
					Category: metrics.CategoryCommandExecution,
					Action:   metrics.ActionApplyApplication,
					Label:    metrics.LabelStart,
				})

				opts.Application, err = commands.InferApplicationFromStdinOrFile(*o.Declaration, o.In, o.FileSystem, opts.File)
				if err != nil {
					return fmt.Errorf("inferring application from stdin or file: %w", err)
				}

				return nil
			},
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("failed validating options: %w", err)
			}

			state := o.StateHandlers(o.StateNodes())

			services, err := o.ClientServices(state)
			if err != nil {
				return fmt.Errorf("acquiring client services: %w", err)
			}

			spin, err := spinner.New("applying application", o.Err)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
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
				reconciliation.NewPostgresReconciler(services.ApplicationPostgresService),
			)

			_, err = scheduler.Run(o.Ctx, state)
			if err != nil {
				return fmt.Errorf("reconciling application: %w", err)
			}

			return commands.WriteApplyApplicationSuccessMessage(commands.WriteApplyApplicationSucessMessageOpts{
				Out:         o.Out,
				Application: opts.Application,
				Cluster:     *o.Declaration,
			})
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			metrics.Publish(metrics.Event{
				Category: metrics.CategoryCommandExecution,
				Action:   metrics.ActionApplyApplication,
				Label:    metrics.LabelEnd,
			})

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "Specify the file path. Use \"-\" for stdin")

	return cmd
}
