package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/context"

	"github.com/logrusorgru/aurora"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller"
	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/cobra"

	"sigs.k8s.io/yaml"
)

type applyClusterOpts struct {
	AWSCredentialsType    string
	GithubCredentialsType string

	File string

	Declaration *v1alpha1.Cluster
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

	var clusterID api.ID

	cmd := &cobra.Command{
		Use:     "cluster -f declaration_file",
		Example: "okctl apply cluster -f cluster.yaml",
		Short:   "apply a cluster definition to the world",
		Long:    "ensures your cluster reflects the declaration of it",
		Args:    cobra.ExactArgs(0),
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) (err error) {
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			o.AWSCredentialsType = opts.AWSCredentialsType
			o.GithubCredentialsType = opts.GithubCredentialsType

			opts.Declaration, err = inferClusterFromStdinOrFile(o.In, opts.File)
			if err != nil {
				return fmt.Errorf("inferring cluster: %w", err)
			}

			err = loadNoUserInputUserData(o, cmd)
			if err != nil {
				return fmt.Errorf("loading application data: %w", err)
			}

			err = loadNoUserInputRepoData(o, opts.Declaration)
			if err != nil {
				return fmt.Errorf("loading repository data: %w", err)
			}

			err = o.InitialiseWithEnvAndAWSAccountID(
				opts.Declaration.Metadata.Environment,
				opts.Declaration.Metadata.AccountID,
			)
			if err != nil {
				return fmt.Errorf("initializing okctl: %w", err)
			}

			err = opts.Declaration.Validate()
			if err != nil {
				return fmt.Errorf("validating cluster declaration: %w", err)
			}

			clusterID = api.ID{
				Region:       opts.Declaration.Metadata.Region,
				AWSAccountID: opts.Declaration.Metadata.AccountID,
				Environment:  opts.Declaration.Metadata.Environment,
				Repository:   o.RepoStateWithEnv.GetMetadata().Name,
				ClusterName:  o.RepoStateWithEnv.GetClusterName(),
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			spin, err := spinner.New("synchronizing", o.Err)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
			}

			services, err := o.ClientServices(spin)
			if err != nil {
				return fmt.Errorf("error getting services: %w", err)
			}

			outputDir, _ := o.GetRepoOutputDir(opts.Declaration.Metadata.Environment)

			desiredTree := controller.CreateDesiredStateTree(opts.Declaration)

			reconciliationManager := reconciler.NewReconcilerManager(&resourcetree.CommonMetadata{
				Ctx:         o.Ctx,
				Out:         o.Out,
				Spin:        spin,
				ClusterID:   clusterID,
				Declaration: opts.Declaration,
			})

			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeALBIngress, reconciler.NewALBIngressReconciler(services.ALBIngressController))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeArgoCD, reconciler.NewArgocdReconciler(services.ArgoCD, services.Github))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeAWSLoadBalancerController, reconciler.NewAWSLoadBalancerControllerReconciler(services.AWSLoadBalancerControllerService))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeAutoscaler, reconciler.NewAutoscalerReconciler(services.Autoscaler))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeBlockstorage, reconciler.NewBlockstorageReconciler(services.Blockstorage))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeCluster, reconciler.NewClusterReconciler(services.Cluster))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeExternalDNS, reconciler.NewExternalDNSReconciler(services.ExternalDNS))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeExternalSecrets, reconciler.NewExternalSecretsReconciler(services.ExternalSecrets))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeGithub, reconciler.NewGithubReconciler(services.Github))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeIdentityManager, reconciler.NewIdentityManagerReconciler(services.IdentityManager))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeVPC, reconciler.NewVPCReconciler(services.Vpc))
			reconciliationManager.AddReconciler(resourcetree.ResourceNodeTypeZone, reconciler.NewZoneReconciler(services.Domain))
			reconciliationManager.AddReconciler(
				resourcetree.ResourceNodeTypeNameserverDelegator,
				reconciler.NewNameserverDelegationReconciler(services.NameserverHandler, services.Domain),
			)

			synchronizeOpts := &controller.SynchronizeOpts{
				ClusterID:               clusterID,
				DesiredTree:             desiredTree,
				ReconciliationManager:   reconciliationManager,
				Fs:                      o.FileSystem,
				OutputDir:               outputDir,
				GithubGetter:            o.RepoStateWithEnv.GetGithub,
				IdentityPoolFetcher:     func() state.IdentityPool { return o.RepoStateWithEnv.GetIdentityPool() },
				CIDRGetter:              func() string { return o.RepoStateWithEnv.GetVPC().CIDR },
				PrimaryHostedZoneGetter: func() *state.HostedZone { return o.RepoStateWithEnv.GetPrimaryHostedZone() },
			}

			err = controller.Synchronize(synchronizeOpts)
			if err != nil {
				return fmt.Errorf("error synchronizing declaration with state: %w", err)
			}

			fmt.Fprintln(o.Out, "\nYour cluster is up to date.")
			fmt.Fprintln(o.Out,
				fmt.Sprintf(
					"\nTo access your cluster, run %s to activate the environment for your cluster",
					aurora.Green(fmt.Sprintf("okctl venv %s", clusterID.Environment)),
				),
			)
			fmt.Fprintln(o.Out, fmt.Sprintf("Your cluster should then be available with %s", aurora.Green("kubectl")))

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.File, "file", "f", "", usageApplyClusterFile)
	flags.StringVarP(&opts.AWSCredentialsType, "aws-credentials-type", "a", context.AWSCredentialsTypeSAML,
		fmt.Sprintf(
			"The form of authentication to use for AWS. Possible values: [%s,%s]",
			context.AWSCredentialsTypeSAML,
			context.AWSCredentialsTypeAccessKey,
		),
	)
	flags.StringVarP(
		&opts.GithubCredentialsType,
		"github-credentials-type",
		"g",
		context.GithubCredentialsTypeDeviceAuthentication,
		fmt.Sprintf(
			"The form of authentication to use for Github. Possible values: [%s,%s]",
			context.GithubCredentialsTypeDeviceAuthentication,
			context.GithubCredentialsTypeToken,
		),
	)

	cmd.Hidden = true

	return cmd
}

const usageApplyClusterFile = `specifies where to read the declaration from. Use "-" for stdin`

func inferClusterFromStdinOrFile(stdin io.Reader, path string) (*v1alpha1.Cluster, error) {
	var (
		inputReader io.Reader
		err         error
	)

	switch path {
	case "-":
		inputReader = stdin
	default:
		inputReader, err = os.Open(filepath.Clean(path))
		if err != nil {
			return nil, fmt.Errorf("unable to read file: %w", err)
		}
	}

	var (
		buffer  bytes.Buffer
		cluster v1alpha1.Cluster
	)

	cluster = v1alpha1.NewDefaultCluster("", "", "", "", "", "")

	_, err = io.Copy(&buffer, inputReader)
	if err != nil {
		return nil, fmt.Errorf("error copying reader data: %w", err)
	}

	err = yaml.Unmarshal(buffer.Bytes(), &cluster)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling buffer: %w", err)
	}

	return &cluster, nil
}

// ~/.okctl.yaml
func loadNoUserInputUserData(o *okctl.Okctl, cmd *cobra.Command) error {
	userDataNotFound := load.CreateOnUserDataNotFoundWithNoInput()

	o.UserDataLoader = load.UserDataFromFlagsEnvConfigDefaults(cmd, userDataNotFound)

	return o.LoadUserData()
}

// .okctl.yaml
func loadNoUserInputRepoData(o *okctl.Okctl, declaration *v1alpha1.Cluster) error {
	repoDataNotFound := load.CreateOnRepoDataNotFoundWithNoUserInput(declaration)

	o.RepoDataLoader = load.RepoDataFromConfigFile(repoDataNotFound)

	return o.LoadRepoData()
}
