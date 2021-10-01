package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/oslokommune/okctl/pkg/upgrade/clusterversion"
	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversion"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/upgrade"

	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/controller/cluster/reconciliation"

	"github.com/asdine/storm/v3/codec/json"

	"github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/logrusorgru/aurora"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/load"
	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

type applyClusterOpts struct {
	DisableSpinner bool
	File           string
	Declaration    *v1alpha1.Cluster
}

// Validate ensures the applyClusterOpts contains the right information
func (o *applyClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.File, validation.Required),
	)
}

// nolint funlen
func buildApplyClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := applyClusterOpts{}
	var originalClusterVersioner originalclusterversion.Versioner

	var clusterVersioner clusterversion.Versioner

	cmd := &cobra.Command{
		Use:     "cluster -f declaration_file",
		Example: "okctl apply cluster -f cluster.yaml",
		Short:   ApplyClusterShortDescription,
		Long:    ApplyClusterLongDescription,
		Args:    cobra.ExactArgs(0),
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) (err error) {
			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				os.Exit(1)
			}()

			enableServiceUserAuthentication(o)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.Declaration, err = commands.InferClusterFromStdinOrFile(o.In, opts.File)
			if err != nil {
				return fmt.Errorf("inferring cluster: %w", err)
			}

			err = opts.Declaration.Validate()
			if err != nil {
				return fmt.Errorf("validating cluster declaration: %w", err)
			}

			err = loadNoUserInputUserData(o, cmd)
			if err != nil {
				return fmt.Errorf("loading application data: %w", err)
			}

			o.Declaration = opts.Declaration

			// Move into a function
			{
				baseDir, err := o.GetRepoDir()
				if err != nil {
					return err
				}

				stormDB := path.Join(baseDir, o.Declaration.Github.OutputPath, o.Declaration.Metadata.Name, constant.DefaultStormDBName)

				exists, err := o.FileSystem.Exists(stormDB)
				if err != nil {
					return err
				}

				if !exists {
					err := o.FileSystem.MkdirAll(path.Dir(stormDB), 0o744)
					if err != nil {
						return err
					}

					db, err := storm.Open(stormDB, storm.Codec(json.Codec))
					if err != nil {
						return err
					}

					err = db.Close()
					if err != nil {
						return err
					}
				}
			}

			err = o.Initialise()
			if err != nil {
				return fmt.Errorf("initializing okctl: %w", err)
			}

			state := o.StateHandlers(o.StateNodes())

			// Cluster version
			clusterID := api.ID{
				Region:       o.Declaration.Metadata.Region,
				AWSAccountID: o.Declaration.Metadata.AccountID,
				ClusterName:  o.Declaration.Metadata.Name,
			}

			clusterVersioner = clusterversion.New(
				o.Out,
				clusterID,
				state.Upgrade,
			)

			err = clusterVersioner.ValidateBinaryVsClusterVersion(version.GetVersionInfo().Version)
			if err != nil {
				return fmt.Errorf(commands.ValidateBinaryVsClusterVersionErr, err)
			}

			// Original version
			originalClusterVersioner = originalclusterversion.New(
				clusterID,
				state.Upgrade,
				state.Cluster,
			)

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			var spinnerWriter io.Writer
			if opts.DisableSpinner {
				spinnerWriter = ioutil.Discard
			} else {
				spinnerWriter = o.Err
			}

			spin, err := spinner.New("applying cluster", spinnerWriter)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
			}

			state := o.StateHandlers(o.StateNodes())

			services, err := o.ClientServices(state)
			if err != nil {
				return fmt.Errorf("error getting services: %w", err)
			}

			schedulerOpts := common.SchedulerOpts{
				Out:                             o.Out,
				Spinner:                         spin,
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
				reconciliation.NewPostgresReconciler(services.Component),
				reconciliation.NewCleanupSGReconciler(o.CloudProvider),
			)

			_, err = scheduler.Run(o.Ctx, state)
			if err != nil {
				return fmt.Errorf("synchronizing declaration with state: %w", err)
			}

			err = handleClusterVersioning(o, originalClusterVersioner, clusterVersioner, opts)
			if err != nil {
				return fmt.Errorf("handle cluster versioning: %w", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.File,
		"file",
		"f",
		"",
		usageApplyClusterFile,
	)
	flags.BoolVar(&opts.DisableSpinner,
		"no-spinner",
		false,
		"disables progress spinner",
	)

	return cmd
}

// (tag UPGR01) In the future, we can replace most of this function with
// originalClusterVersioner.SaveOriginalClusterVersionIfNotExists(version.GetVersionInfo().Version)
// clusterVersioner.SaveClusterVersionIfNotExists(version.GetVersionInfo().Version)
// See those functions' comments.
func handleClusterVersioning(
	o *okctl.Okctl,
	originalClusterVersioner originalclusterversion.Versioner,
	clusterVersioner clusterversion.Versioner,
	opts applyClusterOpts,
) error {
	hasOriginalClusterVersion, err := originalClusterVersioner.OriginalClusterVersionExists()
	if err != nil {
		return fmt.Errorf("checking if original cluster version exists: %w", err)
	}

	err = originalClusterVersioner.SaveOriginalClusterVersionFromClusterTagIfNotExists()
	if err != nil {
		return fmt.Errorf(originalclusterversion.SaveErrorMessage, err)
	}

	err = clusterVersioner.SaveClusterVersionFromOriginalClusterVersionIfNotExists()
	if err != nil {
		return fmt.Errorf(commands.SaveClusterVersionErr, err)
	}

	// When deleting for tag UPGR01, keep this function (or its contents), but delete this comment
	printClusterReadyMessage(o, opts)

	// Remove this when original cluster version has been stored for all users
	clusterVersion, err := clusterVersioner.GetClusterVersion()
	if err != nil {
		return fmt.Errorf("getting original cluster version: %w", err)
	}

	// We show this message only for old (pre upgrade-release) clusters, because for new clusters, we will always store
	// version immediately.
	shouldShowMessage, err := isVersionFromBeforeUpgradeWasReleased(clusterVersion)
	if err != nil {
		return fmt.Errorf("checking upgrade release version: %w", err)
	}

	if !hasOriginalClusterVersion && shouldShowMessage {
		stateFile := path.Join(
			opts.Declaration.Github.OutputPath, opts.Declaration.Metadata.Name, constant.DefaultStormDBName)

		_, _ = fmt.Fprintf(o.Out, "\nOkctl detected that parts of the cluster state had to be "+
			"initialized to support future upgrades. The cluster state has now been initialized. You "+
			"must commit and push changes to %s. For more information, see %s\n",
			stateFile,
			upgrade.DocumentationURL)
	}

	return nil
}

func printClusterReadyMessage(o *okctl.Okctl, opts applyClusterOpts) {
	{
		_, _ = fmt.Fprintln(o.Out, "\nYour cluster is up to date.")
		_, _ = fmt.Fprintf(o.Out,
			"\nTo access your cluster, run %s to activate the environment for your cluster\n",
			aurora.Green(fmt.Sprintf("okctl venv -c %s", opts.File)),
		)
		_, _ = fmt.Fprintf(o.Out, "Your cluster should then be available with %s\n", aurora.Green("kubectl"))
	}
}

func isVersionFromBeforeUpgradeWasReleased(versionString string) (bool, error) {
	v, err := semver.NewVersion(versionString)
	if err != nil {
		return false, fmt.Errorf("cannot create semver version from '%s': %w", versionString, err)
	}

	versionWhereUpgradeWasReleased, err := semver.NewVersion("0.0.67")
	if err != nil {
		return false, fmt.Errorf("cannot create semver version: %w", err)
	}

	return v.LessThan(versionWhereUpgradeWasReleased), nil
}

const usageApplyClusterFile = `specifies where to read the declaration from. Use "-" for stdin`

// ~/.okctl.yaml
func loadNoUserInputUserData(o *okctl.Okctl, cmd *cobra.Command) error {
	userDataNotFound := load.CreateOnUserDataNotFoundWithNoInput()

	o.UserDataLoader = load.UserDataFromFlagsEnvConfigDefaults(cmd, userDataNotFound)

	return o.LoadUserData()
}
