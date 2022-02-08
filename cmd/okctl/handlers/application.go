package handlers

import (
	"fmt"
	"io"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/controller/application/reconciliation"
	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

// HandleApplication knows how to reconcile and purge applications
func HandleApplication(opts *HandleApplicationOpts) RunEHandler {
	return func(cmd *cobra.Command, _ []string) error {
		err := opts.Validate()
		if err != nil {
			return fmt.Errorf("failed validating options: %w", err)
		}

		state := opts.Okctl.StateHandlers(opts.Okctl.StateNodes())

		services, err := opts.Okctl.ClientServices(state)
		if err != nil {
			return fmt.Errorf("acquiring client services: %w", err)
		}

		statusVerb := "applying"
		if opts.Purge {
			statusVerb = "deleting"
		}

		spin, err := spinner.New(fmt.Sprintf("%s application", statusVerb), opts.Okctl.Err)
		if err != nil {
			return fmt.Errorf("error creating spinner: %w", err)
		}

		schedulerOpts := common.SchedulerOpts{
			Out:                             opts.Okctl.Out,
			Spinner:                         spin,
			ReconciliationLoopDelayFunction: common.DefaultDelayFunction,
			ClusterDeclaration:              *opts.Okctl.Declaration,
			ApplicationDeclaration:          opts.Application,
			PurgeFlag:                       opts.Purge,
		}

		scheduler := common.NewScheduler(schedulerOpts,
			reconciliation.NewCertificateReconciler(services.Certificate, services.Domain),
			reconciliation.NewApplicationReconciler(services.ApplicationService, services.ApplicationPostgresService),
			reconciliation.NewContainerRepositoryReconciler(services.ContainerRepository),
			reconciliation.NewPostgresReconciler(services.ApplicationPostgresService),
			reconciliation.NewArgoCDApplicationReconciler(services.ApplicationService),
		)

		_, err = scheduler.Run(opts.Okctl.Ctx, state)
		if err != nil {
			return fmt.Errorf("reconciling application: %w", err)
		}

		return writeSuccessMessage(opts.Okctl.Out, *opts.Okctl.Declaration, opts.Application, opts.Purge)
	}
}

func writeSuccessMessage(out io.Writer, cluster v1alpha1.Cluster, application v1alpha1.Application, purgeFlag bool) error {
	if purgeFlag {
		return commands.WriteDeleteApplicationSuccessMessage(commands.WriteDeleteApplicationSuccessMessageOpts{
			Out:         out,
			Cluster:     cluster,
			Application: application,
		})
	}

	return commands.WriteApplyApplicationSuccessMessage(commands.WriteApplyApplicationSucessMessageOpts{
		Out:         out,
		Cluster:     cluster,
		Application: application,
	})
}

// HandleApplicationOpts contains all the necessary options for application reconciliation
type HandleApplicationOpts struct {
	Okctl *okctl.Okctl

	File string

	Application v1alpha1.Application
	Purge       bool
}

// Validate the options for "apply application"
func (o HandleApplicationOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.File, validation.Required),
		validation.Field(&o.Application),
	)
}
