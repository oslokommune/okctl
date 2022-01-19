package handlers

import (
	"fmt"
	"io"
	"text/template"

	"github.com/AlecAivazis/survey/v2"

	"github.com/oslokommune/okctl/pkg/clients/kubectl/binary"

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
			fmt.Fprintln(opts.Okctl.Err, "user wasn't ready to continue, aborting.")

			return nil
		}

		state := opts.Okctl.StateHandlers(opts.Okctl.StateNodes())

		services, err := opts.Okctl.ClientServices(state)
		if err != nil {
			return fmt.Errorf("acquiring client services: %w", err)
		}

		kubectlClient := binary.New(
			opts.Okctl.FileSystem,
			opts.Okctl.BinariesProvider,
			opts.Okctl.CredentialsProvider,
			*opts.Okctl.Declaration,
		)

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
			reconciliation.NewCertificateReconciler(services.Certificate, services.Domain, kubectlClient),
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

	return deleteApplicationReadyCheck(opts.Okctl.Out, opts.Application, opts.Confirm)
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

// HandleApplicationOpts contains all the necessary options for application reconciliation
type HandleApplicationOpts struct {
	Okctl *okctl.Okctl

	File string

	Application v1alpha1.Application
	Purge       bool
	Confirm     bool
}

// Validate the options for "apply application"
func (o HandleApplicationOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.File, validation.Required),
		validation.Field(&o.Application),
	)
}
