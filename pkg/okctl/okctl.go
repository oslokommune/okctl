// Package okctl implements the core logic for creating providers
// and loading configuration state
package okctl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/oslokommune/okctl/pkg/context"

	"github.com/logrusorgru/aurora/v3"

	clientFilesystem "github.com/oslokommune/okctl/pkg/client/core/store/filesystem"

	"github.com/oslokommune/okctl/pkg/ask"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/core/api/rest"
	"github.com/oslokommune/okctl/pkg/client/core/report/console"
	stateSaver "github.com/oslokommune/okctl/pkg/client/core/state"
	githubClient "github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/spinner"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/oslokommune/okctl/pkg/api/core/store/noop"

	"github.com/oslokommune/okctl/pkg/credentials/github"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api/core"
	awsProvider "github.com/oslokommune/okctl/pkg/api/core/cloudprovider/aws"
	"github.com/oslokommune/okctl/pkg/api/core/run"
	"github.com/oslokommune/okctl/pkg/api/core/store/filesystem"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/credentials/aws/scrape"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/keyring"
	"github.com/oslokommune/okctl/pkg/storage"
)

const authenticatorTimeout = 5 * time.Second

// Okctl stores all state required for invoking commands
type Okctl struct {
	*config.Config

	CloudProvider       v1alpha1.CloudProvider
	BinariesProvider    binaries.Provider
	CredentialsProvider credentials.Provider

	activeEnv  string
	restClient *rest.HTTPClient
}

// InitialiseWithOnlyEnv initialises okctl when the aws account is has been
// set previously
func (o *Okctl) InitialiseWithOnlyEnv(env string) error {
	if !o.RepoState.HasEnvironment(env) {
		return ErrorEnvironmentNotFound{
			TargetEnvironment:     env,
			AvailableEnvironments: getEnvironments(o.RepoState.Clusters),
		}
	}

	repoDir, err := o.GetRepoDir()
	if err != nil {
		return err
	}

	o.RepoStateWithEnv = state.NewRepositoryStateWithEnv(env, o.RepoState, state.DefaultFileSystemSaver(
		config.DefaultRepositoryStateFile,
		repoDir,
		o.FileSystem,
	))

	o.activeEnv = env

	return o.initialise()
}

// InitialiseWithEnvAndAWSAccountID initialises okctl when aws account id hasn't
// been set yet
func (o *Okctl) InitialiseWithEnvAndAWSAccountID(env, awsAccountID string) error {
	repoDir, err := o.GetRepoDir()
	if err != nil {
		return err
	}

	o.RepoStateWithEnv = state.NewRepositoryStateWithEnv(env, o.RepoState, state.DefaultFileSystemSaver(
		config.DefaultRepositoryStateFile,
		repoDir,
		o.FileSystem,
	))

	cluster := o.RepoStateWithEnv.GetCluster()
	cluster.AWSAccountID = awsAccountID
	cluster.Environment = env
	cluster.Name = o.RepoState.Metadata.Name

	_, err = o.RepoStateWithEnv.SaveCluster(cluster)
	if err != nil {
		return err
	}

	o.activeEnv = env

	return o.initialise()
}

// New returns a new okctl instance
func New() *Okctl {
	return &Okctl{
		Config: config.New(),
	}
}

// Binaries returns the application binaries
func (o *Okctl) Binaries() []state.Binary {
	return o.UserState.Binaries
}

// Host returns the host information
func (o *Okctl) Host() state.Host {
	return o.UserState.Host
}

// Username returns the username of the active user
func (o *Okctl) Username() string {
	return o.UserState.User.Username
}

// ErrorFormatter helps add more information to the error message
func (o *Okctl) ErrorFormatter(command, userDir string) func(err error) error {
	const errMsg = `
# Enable debug output
$ OKCTL_DEBUG=true okctl %s
# Inspect the logs
$ cat %s
# Ask for help on slack: %s

%s`

	return func(err error) error {
		return fmt.Errorf(errMsg,
			command,
			path.Join(userDir, config.DefaultLogDir, config.DefaultLogName),
			aurora.Bold("#kjøremiljø-support"),
			err,
		)
	}
}

// ClientServices returns the initialised client-side services
func (o *Okctl) ClientServices(spin spinner.Spinner) (*clientCore.Services, error) {
	outputDir, err := o.GetRepoOutputDir(o.activeEnv)
	if err != nil {
		return nil, err
	}

	applicationsOutputDir, err := o.GetRepoApplicatiosOutputDir()
	if err != nil {
		return nil, err
	}

	ghClient, err := githubClient.New(o.Ctx, o.CredentialsProvider.Github())
	if err != nil {
		return nil, err
	}

	return &clientCore.Services{
		ALBIngressController:             o.albIngressService(outputDir, spin),
		AWSLoadBalancerControllerService: o.awsLoadBalancerControllerService(outputDir, spin),
		ArgoCD:                           o.argocdService(outputDir, spin),
		ApplicationService:               o.applicationService(outputDir, applicationsOutputDir, spin),
		Certificate:                      o.certService(outputDir, spin),
		Cluster:                          o.clusterService(outputDir, spin),
		Domain:                           o.domainService(outputDir, spin),
		ExternalDNS:                      o.externalDNSService(outputDir, spin),
		ExternalSecrets:                  o.externalSecretsService(outputDir, spin),
		Github:                           o.githubService(ghClient, spin),
		Manifest:                         o.manifestService(outputDir, spin),
		NameserverHandler:                o.nameserverHandlerService(ghClient, outputDir, spin),
		Parameter:                        o.paramService(outputDir, spin),
		Vpc:                              o.vpcService(outputDir, spin),
		IdentityManager:                  o.identityManagerService(outputDir, spin),
		Autoscaler:                       o.autoscalerService(outputDir, spin),
		Blockstorage:                     o.blockstorageService(outputDir, spin),
		Monitoring:                       o.monitoringService(outputDir, spin),
	}, nil
}

func (o *Okctl) managedPolicyService(outputDir string, spin spinner.Spinner) client.ManagedPolicyService {
	return clientCore.NewManagedPolicyService(spin,
		rest.NewManagedPolicyAPI(o.restClient),
		clientFilesystem.NewManagedPolicyStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(outputDir, config.DefaultPolicyBaseDir),
			},
			o.FileSystem,
		),
		console.NewManagedPolicyReport(o.Err, spin),
	)
}

func (o *Okctl) serviceAccountService(outputDir string, spin spinner.Spinner) client.ServiceAccountService {
	return clientCore.NewServiceAccountService(
		spin,
		rest.NewServiceAccountAPI(o.restClient),
		clientFilesystem.NewServiceAccountStore(clientFilesystem.Paths{
			OutputFile: config.DefaultServiceAccountOutputsFile,
			ConfigFile: config.DefaultServiceAccountConfigFile,
			BaseDir:    path.Join(outputDir, config.DefaultServiceAccountBaseDir),
		}, o.FileSystem),
		console.NewServiceAccountReport(o.Err, spin),
	)
}

func (o *Okctl) monitoringService(outputDir string, spin spinner.Spinner) client.MonitoringService {
	monitoringDir := path.Join(outputDir, config.DefaultMonitoringBaseDir)

	return clientCore.NewMonitoringService(
		spin,
		rest.NewMonitoringAPI(o.restClient),
		clientFilesystem.NewMonitoringStore(
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, config.DefaultTempoBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, config.DefaultPromtailBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, config.DefaultLokiBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, config.DefaultKubePromStackBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultKubePromStackOutputsFile,
				BaseDir:    path.Join(monitoringDir, config.DefaultKubePromStackBaseDir),
			},
			o.FileSystem,
		),
		stateSaver.NewMonitoringState(o.RepoStateWithEnv),
		console.NewMonitoringReport(o.Err, spin),
		o.certService(monitoringDir, spin.SubSpinner()),
		o.identityManagerService(monitoringDir, spin.SubSpinner()),
		o.manifestService(monitoringDir, spin.SubSpinner()),
		o.paramService(monitoringDir, spin.SubSpinner()),
		o.serviceAccountService(monitoringDir, spin.SubSpinner()),
		o.managedPolicyService(monitoringDir, spin.SubSpinner()),
		o.CloudProvider,
	)
}

func (o *Okctl) identityManagerService(outputDir string, spin spinner.Spinner) client.IdentityManagerService {
	identityPoolBaseDir := path.Join(outputDir, config.DefaultIdentityPoolBaseDir)

	identityManagerService := clientCore.NewIdentityManagerService(
		spin,
		rest.NewIdentityManagerAPI(o.restClient),
		clientFilesystem.NewIdentityManagerStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultIdentityPoolOutputsFile,
				CloudFormationFile: config.DefaultIdentityPoolCloudFormationTemplate,
				BaseDir:            identityPoolBaseDir,
			},
			clientFilesystem.Paths{
				CloudFormationFile: config.DefaultCertificateCloudFormationTemplate,
				BaseDir:            path.Join(identityPoolBaseDir, config.DefaultCertificateBaseDir),
			},
			clientFilesystem.Paths{
				CloudFormationFile: config.DefaultAliasCloudFormationTemplate,
				BaseDir:            path.Join(identityPoolBaseDir, config.DefaultAliasBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:         config.DefaultIdentityPoolClientOutputsFile,
				CloudFormationFile: config.DefaultIdentityPoolClientCloudFormationTemplate,
				BaseDir:            path.Join(identityPoolBaseDir, config.DefaultIdentityPoolClientsBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:         config.DefaultIdentityPoolUserOutputsFile,
				CloudFormationFile: config.DefaultIdentityPoolUserCloudFormationTemplate,
				BaseDir:            path.Join(identityPoolBaseDir, config.DefaultIdentityPoolUsersBaseDir),
			},
			o.FileSystem,
		),
		stateSaver.NewIdentityManagerState(o.RepoStateWithEnv),
		console.NewIdentityManagerReport(o.Err, spin),
		o.certService(outputDir, spin.SubSpinner()),
	)

	return identityManagerService
}

func (o *Okctl) argocdService(outputDir string, spin spinner.Spinner) client.ArgoCDService {
	argoBaseDir := path.Join(outputDir, config.DefaultArgoCDBaseDir)

	argoService := clientCore.NewArgoCDService(
		spin,
		o.identityManagerService(argoBaseDir, spin.SubSpinner()),
		o.certService(argoBaseDir, spin.SubSpinner()),
		o.manifestService(argoBaseDir, spin.SubSpinner()),
		o.paramService(argoBaseDir, spin.SubSpinner()),
		rest.NewArgoCDAPI(o.restClient),
		clientFilesystem.NewArgoCDStore(
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(outputDir, config.DefaultHelmBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultArgoOutputsFile,
				BaseDir:    argoBaseDir,
			},
			o.FileSystem,
		),
		console.NewArgoCDReport(o.Err, spin),
		stateSaver.NewArgoCDState(o.RepoStateWithEnv),
	)

	return argoService
}

func (o *Okctl) paramService(outputDir string, spin spinner.Spinner) client.ParameterService {
	return clientCore.NewParameterService(
		spin,
		rest.NewParameterAPI(o.restClient),
		clientFilesystem.NewParameterStore(
			clientFilesystem.Paths{
				OutputFile: config.DefaultParameterOutputsFile,
				BaseDir:    path.Join(outputDir, config.DefaultParameterBaseDir),
			},
			o.FileSystem,
		),
		console.NewParameterReport(o.Err, spin),
	)
}

func (o *Okctl) manifestService(outputDir string, spin spinner.Spinner) client.ManifestService {
	return clientCore.NewManifestService(
		spin,
		rest.NewManifestAPI(o.restClient),
		clientFilesystem.NewManifestStore(
			clientFilesystem.Paths{
				OutputFile: config.DefaultNamespaceOutputFile,
				ConfigFile: config.DefaultNamespaceConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultNamespaceBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultConfigMapOutputsFile,
				ConfigFile: config.DefaultConfigMapConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultConfigMapBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultStorageClassOutputsFile,
				ConfigFile: config.DefaultStorageClassConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultStorageClassBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultKubeOutputsFile,
				BaseDir:    path.Join(outputDir, config.DefaultExternalSecretsBaseDir),
			},
			o.FileSystem,
		),
		console.NewManifestReport(o.Err, spin),
	)
}

func (o *Okctl) nameserverHandlerService(ghClient githubClient.Githuber, _ string, spin spinner.Spinner) client.NameserverRecordDelegationService {
	return clientCore.NewNameserverHandlerService(
		ghClient,
		spin,
	)
}

func (o *Okctl) certService(outputDir string, spin spinner.Spinner) client.CertificateService {
	return clientCore.NewCertificateService(
		spin,
		rest.NewCertificateAPI(o.restClient),
		clientFilesystem.NewCertificateStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultCertificateOutputsFile,
				CloudFormationFile: config.DefaultCertificateCloudFormationTemplate,
				BaseDir:            path.Join(outputDir, config.DefaultCertificateBaseDir),
			},
			o.FileSystem,
		),
		stateSaver.NewCertificateState(o.RepoStateWithEnv),
		console.NewCertificateReport(o.Err, spin),
	)
}

func (o *Okctl) githubService(ghClient githubClient.Githuber, spin spinner.Spinner) client.GithubService {
	return clientCore.NewGithubService(
		spin,
		rest.NewGithubAPI(
			o.Err,
			ask.New().WithSpinner(spin),
			rest.NewParameterAPI(o.restClient),
			ghClient,
		),
		console.NewGithubReport(o.Err, spin),
		stateSaver.NewGithubState(o.RepoStateWithEnv),
	)
}

func (o *Okctl) vpcService(outputDir string, spin spinner.Spinner) client.VPCService {
	return clientCore.NewVPCService(
		spin,
		rest.NewVPCAPI(o.restClient),
		clientFilesystem.NewVpcStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultVpcOutputs,
				CloudFormationFile: config.DefaultVpcCloudFormationTemplate,
				BaseDir:            path.Join(outputDir, config.DefaultVpcBaseDir),
			},
			o.FileSystem,
		),
		console.NewVPCReport(o.Err, spin),
		stateSaver.NewVpcState(o.RepoStateWithEnv),
	)
}

func (o *Okctl) clusterService(outputDir string, spin spinner.Spinner) client.ClusterService {
	return clientCore.NewClusterService(
		spin,
		rest.NewClusterAPI(o.restClient),
		clientFilesystem.NewClusterStore(
			clientFilesystem.Paths{
				ConfigFile: config.DefaultClusterConfig,
				BaseDir:    path.Join(outputDir, config.DefaultClusterBaseDir),
			},
			o.FileSystem,
		),
		console.NewClusterReport(o.Err, spin),
		stateSaver.NewClusterState(o.RepoStateWithEnv),
	)
}

func (o *Okctl) autoscalerService(outputDir string, spin spinner.Spinner) client.AutoscalerService {
	return clientCore.NewAutoscalerService(
		spin,
		rest.NewAutoscalerAPI(o.restClient),
		clientFilesystem.NewAutoscalerStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(outputDir, config.DefaultAutoscalerBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultAutoscalerBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(outputDir, config.DefaultAutoscalerBaseDir),
			},
			o.FileSystem,
		),
		console.NewAutoscalerReport(o.Err, spin),
	)
}

func (o *Okctl) blockstorageService(outputDir string, spin spinner.Spinner) client.BlockstorageService {
	return clientCore.NewBlockstorageService(
		spin,
		rest.NewBlockstorageAPI(o.restClient),
		clientFilesystem.NewBlockstorageStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(outputDir, config.DefaultBlockstorageBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultBlockstorageBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(outputDir, config.DefaultBlockstorageBaseDir),
			},
			o.FileSystem,
		),
		console.NewBlockstorageReport(o.Err, spin),
		o.manifestService(path.Join(outputDir, config.DefaultBlockstorageBaseDir), spin.SubSpinner()),
	)
}

func (o *Okctl) externalSecretsService(outputDir string, spin spinner.Spinner) client.ExternalSecretsService {
	return clientCore.NewExternalSecretsService(
		spin,
		rest.NewExternalSecretsAPI(o.restClient),
		clientFilesystem.NewExternalSecretsStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(outputDir, config.DefaultExternalSecretsBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultExternalSecretsBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(outputDir, config.DefaultExternalSecretsBaseDir),
			},
			o.FileSystem,
		),
		console.NewExternalSecretsReport(o.Err, spin),
	)
}

func (o *Okctl) albIngressService(outputDir string, spin spinner.Spinner) client.ALBIngressControllerService {
	return clientCore.NewALBIngressControllerService(
		spin,
		rest.NewALBIngressControllerAPI(o.restClient),
		clientFilesystem.NewALBIngressControllerStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(outputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(outputDir, config.DefaultAlbIngressControllerBaseDir),
			},
			o.FileSystem,
		),
		console.NewAlbIngressControllerReport(o.Err, spin),
	)
}

func (o *Okctl) awsLoadBalancerControllerService(outputDir string, spin spinner.Spinner) client.AWSLoadBalancerControllerService {
	return clientCore.NewAWSLoadBalancerControllerService(
		spin,
		rest.NewAWSLoadBalancerControllerAPI(o.restClient),
		clientFilesystem.NewAWSLoadBalancerControllerStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(outputDir, config.DefaultAWSLoadBalancerControllerBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultAWSLoadBalancerControllerBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  config.DefaultHelmOutputsFile,
				ReleaseFile: config.DefaultHelmReleaseFile,
				ChartFile:   config.DefaultHelmChartFile,
				BaseDir:     path.Join(outputDir, config.DefaultAWSLoadBalancerControllerBaseDir),
			},
			o.FileSystem,
		),
		console.NewAWSLoadBalancerControllerReport(o.Err, spin),
	)
}

func (o *Okctl) applicationService(infrastructureOutputDir, applicationOutputDir string, spin spinner.Spinner) client.ApplicationService {
	return clientCore.NewApplicationService(
		o.FileSystem,
		spin,
		clientFilesystem.Paths{
			BaseDir: applicationOutputDir,
		},
		o.certService(path.Join(infrastructureOutputDir, config.DefaultCertificateBaseDir), spin.SubSpinner()),
		clientFilesystem.NewApplicationStore(
			clientFilesystem.Paths{
				BaseDir: applicationOutputDir,
			},
			o.FileSystem,
		),
		console.NewApplicationReport(o.Out, spin),
	)
}

func (o *Okctl) domainService(outputDir string, spin spinner.Spinner) client.DomainService {
	return clientCore.NewDomainService(
		spin,
		o.Err,
		ask.New().WithSpinner(spin),
		rest.NewDomainAPI(o.restClient),
		clientFilesystem.NewDomainStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultDomainOutputsFile,
				CloudFormationFile: config.DefaultDomainCloudFormationTemplate,
				BaseDir:            path.Join(outputDir, config.DefaultDomainBaseDir),
			},
			o.FileSystem,
		),
		console.NewDomainReport(o.Err, spin),
		stateSaver.NewDomainState(o.RepoStateWithEnv),
	)
}

func (o *Okctl) externalDNSService(outputDir string, spin spinner.Spinner) client.ExternalDNSService {
	return clientCore.NewExternalDNSService(
		spin,
		rest.NewExternalDNSAPI(o.restClient),
		clientFilesystem.NewExternalDNSStore(
			clientFilesystem.Paths{
				OutputFile:         config.DefaultPolicyOutputFile,
				CloudFormationFile: config.DefaultPolicyCloudFormationTemplateFile,
				BaseDir:            path.Join(outputDir, config.DefaultExternalDNSBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultServiceAccountOutputsFile,
				ConfigFile: config.DefaultServiceAccountConfigFile,
				BaseDir:    path.Join(outputDir, config.DefaultExternalDNSBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: config.DefaultKubeOutputsFile,
				BaseDir:    path.Join(outputDir, config.DefaultExternalDNSBaseDir),
			},
			o.FileSystem,
		),
		console.NewExternalDNSReport(o.Err, spin),
	)
}

// Initialise okctl for receiving requests
// nolint: funlen
func (o *Okctl) initialise() error {
	err := o.EnableFileLog()
	if err != nil {
		return err
	}

	err = o.initialiseProviders()
	if err != nil {
		return err
	}

	o.restClient = rest.New(o.Debug, ioutil.Discard, o.ServerURL)

	homeDir, err := o.GetHomeDir()
	if err != nil {
		return err
	}

	appDir, err := o.GetUserDataDir()
	if err != nil {
		return err
	}

	outputDir, err := o.GetRepoOutputDir(o.activeEnv)
	if err != nil {
		return err
	}

	clusterName := o.RepoStateWithEnv.GetClusterName()

	kubeConfigStore := filesystem.NewKubeConfigStore(
		o.CloudProvider,
		config.DefaultClusterKubeConfig,
		path.Join(appDir, config.DefaultCredentialsDirName, clusterName),
		config.DefaultClusterConfig,
		path.Join(outputDir, config.DefaultClusterBaseDir),
		o.FileSystem,
	)

	vpcStore := noop.NewVpcStore()
	clusterStore := noop.NewClusterStore()
	helmStore := noop.NewHelmStore()
	kubeStore := noop.NewKubeStore()
	certificateStore := noop.NewCertificateStore()
	parameterStore := noop.NewParameterStore()

	vpcService := core.NewVpcService(
		awsProvider.NewVpcCloud(o.CloudProvider),
		vpcStore,
	)

	clusterService := core.NewClusterService(
		clusterStore,
		run.NewClusterRun(
			o.Debug,
			kubeConfigStore,
			path.Join(appDir, config.DefaultCredentialsDirName, clusterName, config.DefaultClusterAwsConfig),
			path.Join(appDir, config.DefaultCredentialsDirName, clusterName, config.DefaultClusterAwsCredentials),
			o.BinariesProvider,
			o.CloudProvider,
		),
	)

	managedPolicyService := core.NewManagedPolicyService(awsProvider.NewManagedPolicyCloudProvider(o.CloudProvider))

	serviceAccountService := core.NewServiceAccountService(
		run.NewServiceAccountRun(
			o.Debug,
			path.Join(appDir, config.DefaultCredentialsDirName, clusterName, config.DefaultClusterAwsConfig),
			path.Join(appDir, config.DefaultCredentialsDirName, clusterName, config.DefaultClusterAwsCredentials),
			o.BinariesProvider,
		),
	)

	kubeService := core.NewKubeService(
		kubeStore,
		run.NewKubeRun(o.CloudProvider, o.CredentialsProvider.Aws()),
	)

	awsIamAuth, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return err
	}

	helmRun := run.NewHelmRun(
		helm.New(&helm.Config{
			HomeDir:              homeDir,
			Path:                 fmt.Sprintf("/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:%s", path.Dir(awsIamAuth.BinaryPath)),
			HelmPluginsDirectory: path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
			HelmRegistryConfig:   path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
			HelmRepositoryConfig: path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
			HelmRepositoryCache:  path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
			HelmBaseDir:          path.Join(appDir, config.DefaultHelmBaseDir),
			Debug:                o.Debug,
			DebugOutput:          o.Err,
		},
			o.CredentialsProvider.Aws(),
			o.FileSystem,
		),
		kubeConfigStore,
	)

	helmService := core.NewHelmService(
		helmRun,
		helmStore,
	)

	domainService := core.NewDomainService(
		awsProvider.NewDomainCloudProvider(o.CloudProvider),
	)

	certificateService := core.NewCertificateService(
		awsProvider.NewCertificateCloudProvider(o.CloudProvider),
		certificateStore,
	)

	parameterService := core.NewParameterService(
		awsProvider.NewParameterCloudProvider(o.CloudProvider),
		parameterStore,
	)

	// When creating a certificate for a CloudFront distribution, we
	// need to create the certificate in us-east-1
	provider, err := o.NewCloudProviderWithRegion("us-east-1")
	if err != nil {
		return err
	}

	identityManagerService := core.NewIdentityManagerService(
		awsProvider.NewIdentityManagerCloudProvider(o.CloudProvider),
		awsProvider.NewCertificateCloudProvider(provider),
	)

	services := core.Services{
		Cluster:         clusterService,
		Vpc:             vpcService,
		ManagedPolicy:   managedPolicyService,
		ServiceAccount:  serviceAccountService,
		Helm:            helmService,
		Kube:            kubeService,
		Domain:          domainService,
		Certificate:     certificateService,
		Parameter:       parameterService,
		IdentityManager: identityManagerService,
	}

	endpoints := core.GenerateEndpoints(services, core.InstrumentEndpoints(o.Logger))

	handlers := core.MakeHandlers(core.EncodeJSONResponse, endpoints)

	router := http.NewServeMux()
	router.Handle("/", core.AttachRoutes(handlers))

	server := &http.Server{
		Handler: router,
		Addr:    o.Destination,
	}

	// nolint: gomnd
	errs := make(chan error, 2)

	go func() {
		errs <- server.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return nil
}

// initialiseProviders knows how to create all required providers
func (o *Okctl) initialiseProviders() error {
	err := o.newCredentialsProvider()
	if err != nil {
		return err
	}

	err = o.newBinariesProvider()
	if err != nil {
		return err
	}

	err = o.newCloudProvider()
	if err != nil {
		return err
	}

	return nil
}

// newBinariesProvider creates a provider for loading binaries
func (o *Okctl) newBinariesProvider() error {
	userDataDir, err := o.GetUserDataDir()
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(
		o.Err,
		o.Logger,
		true,
		o.Host(),
		o.Binaries(),
		storage.NewFileSystemStorage(userDataDir),
	)
	if err != nil {
		return errors.E(err, "failed to create binaries fetcher", errors.Internal)
	}

	out := ioutil.Discard
	if o.Debug {
		out = o.Err
	}

	o.BinariesProvider = binaries.New(o.Logger, out, o.CredentialsProvider.Aws(), fetcher)

	return nil
}

// NewCloudProviderWithRegion create a cloud provider with a spesific region
func (o *Okctl) NewCloudProviderWithRegion(region string) (v1alpha1.CloudProvider, error) {
	c, err := cloud.New(region, o.CredentialsProvider.Aws())
	if err != nil {
		return nil, err
	}

	return c.Provider, nil
}

// newCloudProvider creates a provider for running cloud operations
func (o *Okctl) newCloudProvider() error {
	c, err := cloud.New(o.RepoStateWithEnv.GetMetadata().Region, o.CredentialsProvider.Aws())
	if err != nil {
		return err
	}

	o.CloudProvider = c.Provider

	return nil
}

func (o *Okctl) getAWSAuthenticator() (*aws.Auth, error) {
	if o.AWSCredentialsType == context.AWSCredentialsTypeAccessKey {
		return aws.New(aws.NewInMemoryStorage(), aws.NewAuthEnvironment(o.RepoStateWithEnv.GetMetadata().Region, os.Getenv)), nil
	}

	appDir, err := o.GetUserDataDir()
	if err != nil {
		return nil, err
	}

	defaultRing, err := keyring.DefaultKeyringForOS()
	if err != nil {
		return nil, fmt.Errorf(`unable to create a keyring. It is possible no valid backends were found 
on your system, take a look at this site for valid options:
https://github.com/99designs/keyring#keyring

On linux pass works well, for instance:
https://www.passwordstore.org/

%w`, err)
	}

	k, err := keyring.New(defaultRing, o.Debug)
	if err != nil {
		return nil, err
	}

	storedPassword, _ := k.Fetch(keyring.KeyTypeUserPassword)
	fn := func(username, password string) {
		// We do not handle this error, since we do not want the process to stop even if we cannot
		// save password
		_ = k.Store(keyring.KeyTypeUserPassword, password)
	}

	authStore := aws.NewIniPersister(aws.NewFileSystemIniStorer(
		config.DefaultClusterAwsConfig,
		config.DefaultClusterAwsCredentials,
		path.Join(appDir, config.DefaultCredentialsDirName, o.RepoStateWithEnv.GetClusterName()),
		o.FileSystem,
	))

	return aws.New(authStore, aws.NewAuthSAML(
		o.RepoStateWithEnv.GetCluster().AWSAccountID,
		o.RepoStateWithEnv.GetMetadata().Region,
		scrape.New(),
		aws.DefaultStsProvider,
		aws.Interactive(o.Username(), storedPassword, fn),
	)), nil
}

func (o *Okctl) getGithubAuthenticator() (*github.Auth, error) {
	if o.GithubCredentialsType == context.GithubCredentialsTypeToken {
		return github.New(
			github.NewInMemoryPersister(),
			&http.Client{Timeout: authenticatorTimeout},
			github.NewAuthEnvironment(os.Getenv),
		), nil
	}

	defaultRing, err := keyring.DefaultKeyringForOS()
	if err != nil {
		return nil, fmt.Errorf(`unable to create a keyring. It is possible no valid backends were found 
on your system, take a look at this site for valid options:
https://github.com/99designs/keyring#keyring

On linux pass works well, for instance:
https://www.passwordstore.org/

%w`, err)
	}

	k, err := keyring.New(defaultRing, o.Debug)
	if err != nil {
		return nil, err
	}

	return github.New(
		github.NewKeyringPersister(k),
		&http.Client{Timeout: authenticatorTimeout},
		github.NewAuthDeviceFlow(github.DefaultGithubOauthClientID, github.RequiredScopes()),
	), nil
}

// newCredentialsProvider knows how to load credentials
func (o *Okctl) newCredentialsProvider() error {
	awsAuthenticator, err := o.getAWSAuthenticator()
	if err != nil {
		return fmt.Errorf("error acquiring AWS authenticator: %w", err)
	}

	githubAuthenticator, err := o.getGithubAuthenticator()
	if err != nil {
		return fmt.Errorf("error acquiring Github authenticator: %w", err)
	}

	o.CredentialsProvider = credentials.New(awsAuthenticator, githubAuthenticator)

	return nil
}
