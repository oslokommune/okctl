package handlers

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/oslokommune/okctl/pkg/client/core"

	"github.com/AlecAivazis/survey/v2"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/controller/application/reconciliation"
	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

// HandleApplication knows how to reconcile and purge applications
func HandleApplication(opts *HandleApplicationOpts) RunEHandler { //nolint:funlen
	return func(cmd *cobra.Command, _ []string) error {
		err := opts.Validate()
		if err != nil {
			return fmt.Errorf("failed validating options: %w", err)
		}

		ready, err := reconcileApplicationReadyCheck(opts)
		if err != nil {
			return fmt.Errorf("prompting user: %w", err)
		}

		if !ready {
			fmt.Fprintln(opts.Err, "user wasn't ready to continue, aborting.")

			return nil
		}

		statusVerb := "applying"
		if opts.Purge {
			statusVerb = "deleting"
		}

		spin, err := spinner.New(fmt.Sprintf("%s application", statusVerb), opts.Err)
		if err != nil {
			return fmt.Errorf("error creating spinner: %w", err)
		}

		scheduler := CreateScheduler(CreateSchedulerOpts{
			Out:                 opts.Out,
			Services:            opts.Services,
			State:               opts.State,
			Spinner:             spin,
			ClusterManifest:     opts.ClusterManifest,
			ApplicationManifest: opts.ApplicationManifest,
			Purge:               opts.Purge,
			DelayFunction:       opts.DelayFunction,
		})

		_, err = scheduler.Run(opts.Ctx, opts.State)
		if err != nil {
			return fmt.Errorf("reconciling application: %w", err)
		}

		return writeSuccessMessage(opts.Out, opts.ClusterManifest, opts.ApplicationManifest, opts.Purge)
	}
}

// CreateScheduler knows how to create a complete application reconciliation scheduler
func CreateScheduler(opts CreateSchedulerOpts) common.Scheduler {
	schedulerOpts := common.SchedulerOpts{
		Out:                             opts.Out,
		Spinner:                         opts.Spinner,
		ReconciliationLoopDelayFunction: opts.DelayFunction,
		ClusterDeclaration:              opts.ClusterManifest,
		ApplicationDeclaration:          opts.ApplicationManifest,
		PurgeFlag:                       opts.Purge,
	}

	scheduler := common.NewScheduler(schedulerOpts,
		reconciliation.NewCertificateReconciler(opts.Services.Certificate, opts.Services.Domain),
		reconciliation.NewApplicationReconciler(opts.Services.ApplicationService, opts.Services.ApplicationPostgresService),
		reconciliation.NewContainerRepositoryReconciler(opts.Services.ContainerRepository),
		reconciliation.NewPostgresReconciler(opts.Services.ApplicationPostgresService),
		reconciliation.NewArgoCDApplicationReconciler(opts.Services.ApplicationService),
	)

	return scheduler
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

const deleteApplicationPromptTemplate = `
This will delete the application {{ .ApplicationName }} and all associated resources. Resources include:
{{- if .HasIngress }}{{ printf "\n" }}- Amazon Load Balancer instance{{ end }}
{{- if .HasIngress }}{{ printf "\n" }}- SSL certificate{{ end }}
{{- if .HasECR }}{{ printf "\n" }}- ECR repository and all Docker images inside{{ end }} 
- ArgoCD application manifest locally and in the remote IAC repository
- Kubernetes manifests locally

`

type deleteApplicationPromptTemplateOpts struct {
	ApplicationName string
	HasIngress      bool
	HasECR          bool
}

func writeDeleteApplicationReadyCheckInfo(out io.Writer, opts deleteApplicationPromptTemplateOpts) error {
	t, err := template.New("").Parse(deleteApplicationPromptTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(out, opts)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}

func reconcileApplicationReadyCheck(opts *HandleApplicationOpts) (bool, error) {
	if !opts.Purge {
		return true, nil
	}

	return deleteApplicationReadyCheck(opts.Out, opts.ApplicationManifest, opts.Confirm)
}

func deleteApplicationReadyCheck(out io.Writer, application v1alpha1.Application, preConfirmed bool) (bool, error) {
	if preConfirmed {
		return true, nil
	}

	err := writeDeleteApplicationReadyCheckInfo(out, deleteApplicationPromptTemplateOpts{
		ApplicationName: application.Metadata.Name,
		HasIngress:      application.HasIngress(),
		HasECR:          application.Image.HasName(),
	})
	if err != nil {
		return false, fmt.Errorf("printing delete detailed info: %w", err)
	}

	ready := false
	prompt := &survey.Confirm{Message: "are you sure you want to continue?"}

	err = survey.AskOne(prompt, &ready)
	if err != nil {
		return false, fmt.Errorf("prompting user: %w", err)
	}

	return ready, nil
}

// CreateSchedulerOpts defines necessary data to create an application reconciliation scheduler
type CreateSchedulerOpts struct {
	Out io.Writer

	Services *core.Services
	State    *core.StateHandlers

	Spinner             spinner.Spinner
	ClusterManifest     v1alpha1.Cluster
	ApplicationManifest v1alpha1.Application
	Purge               bool
	DelayFunction       func()
}

// HandleApplicationOpts contains all the necessary options for application reconciliation
type HandleApplicationOpts struct {
	Out io.Writer
	Err io.Writer
	Ctx context.Context

	State    *core.StateHandlers
	Services *core.Services
	File     string

	ClusterManifest     v1alpha1.Cluster
	ApplicationManifest v1alpha1.Application
	Purge               bool
	Confirm             bool
	DelayFunction       func()
}

// Validate the options for "apply application"
func (o HandleApplicationOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.File, validation.Required),
		validation.Field(&o.ApplicationManifest),
	)
}
