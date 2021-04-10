package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/asdine/storm/v3/codec/json"

	"github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/servicequota"

	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/context"

	"github.com/logrusorgru/aurora"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/controller"
	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
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

			vpc, err := o.StateHandlers(o.StateNodes()).
				Vpc.GetVpc(cfn.NewStackNamer().Vpc(o.Declaration.Metadata.Name))
			if err != nil && !errors.Is(err, storm.ErrNotFound) {
				return err
			}

			vpcProvisioned := vpc != nil

			err = servicequota.CheckQuotas(
				servicequota.NewVpcCheck(vpcProvisioned, constant.DefaultRequiredVpcs, o.CloudProvider),
				servicequota.NewEipCheck(vpcProvisioned, constant.DefaultRequiredEpis, o.CloudProvider),
				servicequota.NewIgwCheck(vpcProvisioned, constant.DefaultRequiredIgws, o.CloudProvider),
			)
			if err != nil {
				return fmt.Errorf("checking service quotas: %w", err)
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

			outputDir, _ := o.GetRepoOutputDir()

			reconciliationManager := reconciler.NewCompositeReconciler(spin,
				reconciler.NewArgocdReconciler(services.ArgoCD, services.Github),
				reconciler.NewAWSLoadBalancerControllerReconciler(services.AWSLoadBalancerControllerService),
				reconciler.NewAutoscalerReconciler(services.Autoscaler),
				reconciler.NewKubePrometheusStackReconciler(services.Monitoring),
				reconciler.NewLokiReconciler(services.Monitoring),
				reconciler.NewPromtailReconciler(services.Monitoring),
				reconciler.NewTempoReconciler(services.Monitoring),
				reconciler.NewBlockstorageReconciler(services.Blockstorage),
				reconciler.NewClusterReconciler(services.Cluster),
				reconciler.NewExternalDNSReconciler(services.ExternalDNS),
				reconciler.NewExternalSecretsReconciler(services.ExternalSecrets),
				reconciler.NewIdentityManagerReconciler(services.IdentityManager),
				reconciler.NewVPCReconciler(services.Vpc),
				reconciler.NewZoneReconciler(services.Domain),
				reconciler.NewNameserverDelegationReconciler(services.NameserverHandler),
				reconciler.NewNameserverDelegatedTestReconciler(services.Domain),
				reconciler.NewUsersReconciler(services.IdentityManager),
				reconciler.NewPostgresReconciler(services.Component),
			)

			reconciliationManager.SetCommonMetadata(&resourcetree.CommonMetadata{
				Ctx:         o.Ctx,
				Out:         o.Out,
				ClusterID:   id,
				Declaration: opts.Declaration,
			})

			synchronizeOpts := &controller.SynchronizeOpts{
				Debug:                 o.Debug,
				Out:                   o.Out,
				ID:                    id,
				ClusterDeclaration:    opts.Declaration,
				ReconciliationManager: reconciliationManager,
				Fs:                    o.FileSystem,
				OutputDir:             outputDir,
				StateHandlers:         o.StateHandlers(o.StateNodes()),
			}

			err = controller.Synchronize(synchronizeOpts)
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
