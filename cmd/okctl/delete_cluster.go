package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/api/core/run"
	"github.com/oslokommune/okctl/pkg/controller/cluster/reconciliation"
	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/metrics"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

// DeleteClusterOpts contains the required inputs
type DeleteClusterOpts struct {
	ClusterDeclarationPath string
	DisableSpinner         bool
	Confirm                bool
}

// nolint: gocyclo, funlen, gocognit
func buildDeleteClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &DeleteClusterOpts{}

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: DeleteClusterShortDescription,
		Long:  DeleteClusterLongDescription,
		Args:  cobra.ExactArgs(deleteClusterArgs),
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionDeleteCluster),
			hooks.LoadClusterDeclaration(o, &opts.ClusterDeclarationPath),
			hooks.InitializeOkctl(o),
			hooks.AcquireStateLock(o),
			hooks.DownloadState(o, true),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			var spinnerWriter io.Writer
			if opts.DisableSpinner {
				spinnerWriter = ioutil.Discard
			} else {
				spinnerWriter = o.Err
			}

			spin, err := spinner.New("deleting cluster", spinnerWriter)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
			}

			state := o.StateHandlers(o.StateNodes())

			services, err := o.ClientServices(state)
			if err != nil {
				return fmt.Errorf("error getting services: %w", err)
			}

			kubeService := core.NewKubeService(run.NewKubeRun(o.CloudProvider, o.CredentialsProvider.Aws()))

			schedulerOpts := common.SchedulerOpts{
				Out:                             o.Out,
				Spinner:                         spin,
				PurgeFlag:                       true,
				ReconciliationLoopDelayFunction: common.DefaultDelayFunction,
				ClusterDeclaration:              *o.Declaration,
			}

			scheduler := common.NewScheduler(schedulerOpts,
				reconciliation.NewZoneReconciler(services.Domain),
				reconciliation.NewVPCReconciler(services.Vpc, o.CloudProvider),
				reconciliation.NewNameserverDelegationReconciler(services.NameserverHandler),
				reconciliation.NewClusterReconciler(services.Cluster, o.CloudProvider),
				reconciliation.NewAutoscalerReconciler(services.Autoscaler),
				reconciliation.NewAWSLoadBalancerControllerReconciler(services.AWSLoadBalancerControllerService),
				reconciliation.NewBlockstorageReconciler(services.Blockstorage),
				reconciliation.NewExternalDNSReconciler(services.ExternalDNS),
				reconciliation.NewExternalSecretsReconciler(services.ExternalSecrets),
				reconciliation.NewNameserverDelegatedTestReconciler(services.Domain),
				reconciliation.NewIdentityManagerReconciler(services.IdentityManager),
				reconciliation.NewArgocdReconciler(services.ArgoCD, services.Github),
				reconciliation.NewLokiReconciler(services.Monitoring),
				reconciliation.NewPromtailReconciler(services.Monitoring),
				reconciliation.NewTempoReconciler(services.Monitoring),
				reconciliation.NewKubePrometheusStackReconciler(services.Monitoring),
				reconciliation.NewUsersReconciler(services.IdentityManager),
				reconciliation.NewPostgresReconciler(kubeService, services.Component),
				reconciliation.NewCleanupSGReconciler(o.CloudProvider),
			)

			ready, err := checkIfReady(o.Declaration.Metadata.Name, o, opts.Confirm)
			if err != nil {
				return fmt.Errorf("checking if user is ready: %w", err)
			}

			if !ready {
				return hooks.ReleaseStateLock(o)(nil, nil)
			}

			_, err = scheduler.Run(o.Ctx, state)
			if err != nil {
				return fmt.Errorf("synchronizing declaration with state: %w", err)
			}

			err = hooks.PurgeRemoteState(o)(cmd, nil)
			if err != nil {
				return fmt.Errorf("purging remote state: %w", err)
			}

			return nil
		},
		PostRunE: hooks.RunECombinator(
			hooks.ClearLocalState(o),
			hooks.EmitEndCommandExecutionEvent(metrics.ActionDeleteCluster),
		),
	}
	addAuthenticationFlags(cmd)

	flags := cmd.Flags()
	flags.StringVarP(&opts.ClusterDeclarationPath,
		"file",
		"f",
		"",
		usageApplyClusterFile,
	)
	flags.BoolVar(
		&opts.DisableSpinner,
		"no-spinner",
		false,
		"disables progress spinner",
	)
	flags.BoolVarP(
		&opts.Confirm,
		"confirm",
		"y",
		false,
		"confirm all choices",
	)

	return cmd
}

func checkIfReady(clusterName string, o *okctl.Okctl, preConfirmed bool) (bool, error) {
	if preConfirmed {
		return true, nil
	}

	ready := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("This will delete %s and all assosicated resources, are you sure you want to continue?", clusterName),
	}

	err := survey.AskOne(prompt, &ready)
	if err != nil {
		return false, err
	}

	if !ready {
		_, err = fmt.Fprintf(o.Err, "user wasn't ready to continue, aborting.")
		if err != nil {
			return false, err
		}

		return false, err
	}

	return true, nil
}
