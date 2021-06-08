package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/oslokommune/okctl/pkg/controller/cluster"

	"github.com/oslokommune/okctl/pkg/controller/cluster/reconciliation"

	"github.com/asdine/storm/v3/codec/json"

	"github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/context"

	"github.com/logrusorgru/aurora"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/load"
	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/controller/common/resourcetree"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"
)

type applyClusterOpts struct {
	AWSCredentialsType    string
	GithubCredentialsType string
	DisableSpinner        bool
	File                  string
	Declaration           *v1alpha1.Cluster
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

	var id api.ID

	cmd := &cobra.Command{
		Use:     "cluster -f declaration_file",
		Example: "okctl apply cluster -f cluster.yaml",
		Short:   "apply a cluster definition to the world",
		Long:    "ensures your cluster reflects the declaration of it",
		Args:    cobra.ExactArgs(0),
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) (err error) {
			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				os.Exit(1)
			}()

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			o.AWSCredentialsType = opts.AWSCredentialsType
			o.GithubCredentialsType = opts.GithubCredentialsType

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

			id = api.ID{
				Region:       opts.Declaration.Metadata.Region,
				AWSAccountID: opts.Declaration.Metadata.AccountID,
				ClusterName:  opts.Declaration.Metadata.Name,
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			var spinnerWriter io.Writer
			if opts.DisableSpinner {
				spinnerWriter = ioutil.Discard
			} else {
				spinnerWriter = o.Err
			}

			spin, err := spinner.New("synchronizing", spinnerWriter)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
			}

			handlers := o.StateHandlers(o.StateNodes())

			services, err := o.ClientServices(handlers)
			if err != nil {
				return fmt.Errorf("error getting services: %w", err)
			}

			reconciliationManager := common.NewCompositeReconciler(spin,
				reconciliation.NewArgocdReconciler(services.ArgoCD, services.Github),
				reconciliation.NewAWSLoadBalancerControllerReconciler(services.AWSLoadBalancerControllerService),
				reconciliation.NewAutoscalerReconciler(services.Autoscaler),
				reconciliation.NewKubePrometheusStackReconciler(services.Monitoring),
				reconciliation.NewLokiReconciler(services.Monitoring),
				reconciliation.NewPromtailReconciler(services.Monitoring),
				reconciliation.NewTempoReconciler(services.Monitoring),
				reconciliation.NewBlockstorageReconciler(services.Blockstorage),
				reconciliation.NewClusterReconciler(services.Cluster),
				reconciliation.NewExternalDNSReconciler(services.ExternalDNS),
				reconciliation.NewExternalSecretsReconciler(services.ExternalSecrets),
				reconciliation.NewIdentityManagerReconciler(services.IdentityManager),
				reconciliation.NewVPCReconciler(services.Vpc),
				reconciliation.NewZoneReconciler(services.Domain),
				reconciliation.NewNameserverDelegationReconciler(services.NameserverHandler),
				reconciliation.NewNameserverDelegatedTestReconciler(services.Domain),
				reconciliation.NewUsersReconciler(services.IdentityManager),
				reconciliation.NewPostgresReconciler(services.Component),
				reconciliation.NewCleanupALBReconciler(o.CloudProvider),
				reconciliation.NewCleanupSGReconciler(o.CloudProvider),
				&reconciliation.PostgresGroupReconciler{},
				reconciliation.NewServiceQuotaReconciler(o.CloudProvider),
			)

			reconciliationManager.SetCommonMetadata(&resourcetree.CommonMetadata{
				Ctx:         o.Ctx,
				Out:         o.Out,
				ClusterID:   id,
				Declaration: opts.Declaration,
			})

			synchronizeOpts := &cluster.SynchronizeOpts{
				Debug:                 o.Debug,
				Out:                   o.Out,
				ID:                    id,
				ClusterDeclaration:    opts.Declaration,
				ReconciliationManager: reconciliationManager,
				State:                 handlers,
			}

			err = cluster.Synchronize(synchronizeOpts)
			if err != nil {
				return fmt.Errorf("synchronizing declaration with state: %w", err)
			}

			_, _ = fmt.Fprintln(o.Out, "\nYour cluster is up to date.")
			_, _ = fmt.Fprintln(o.Out,
				fmt.Sprintf(
					"\nTo access your cluster, run %s to activate the environment for your cluster",
					aurora.Green(fmt.Sprintf("okctl venv %s", id.ClusterName)),
				),
			)
			_, _ = fmt.Fprintln(o.Out, fmt.Sprintf("Your cluster should then be available with %s", aurora.Green("kubectl")))

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.AWSCredentialsType,
		"aws-credentials-type",
		"a",
		context.AWSCredentialsTypeSAML,
		fmt.Sprintf(
			"The form of authentication to use for AWS. Possible values: [%s,%s]",
			context.AWSCredentialsTypeSAML,
			context.AWSCredentialsTypeAccessKey,
		),
	)
	flags.StringVarP(&opts.File,
		"file",
		"f",
		"",
		usageApplyClusterFile,
	)
	flags.StringVarP(&opts.GithubCredentialsType,
		"github-credentials-type",
		"g",
		context.GithubCredentialsTypeDeviceAuthentication,
		fmt.Sprintf(
			"The form of authentication to use for Github. Possible values: [%s,%s]",
			context.GithubCredentialsTypeDeviceAuthentication,
			context.GithubCredentialsTypeToken,
		),
	)
	flags.BoolVar(&opts.DisableSpinner,
		"no-spinner",
		false,
		"disables progress spinner",
	)

	return cmd
}

const usageApplyClusterFile = `specifies where to read the declaration from. Use "-" for stdin`

// ~/.okctl.yaml
func loadNoUserInputUserData(o *okctl.Okctl, cmd *cobra.Command) error {
	userDataNotFound := load.CreateOnUserDataNotFoundWithNoInput()

	o.UserDataLoader = load.UserDataFromFlagsEnvConfigDefaults(cmd, userDataNotFound)

	return o.LoadUserData()
}
