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

	"github.com/asdine/storm/v3/codec/json"

	stormpkg "github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/client/core/state/storm"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/constant"

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

	StormDB *stormpkg.DB

	activeEnv       string
	restClient      *rest.HTTPClient
	kubeConfigStore api.KubeConfigStore
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
		constant.DefaultRepositoryStateFile,
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
		constant.DefaultRepositoryStateFile,
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
			path.Join(userDir, constant.DefaultLogDir, constant.DefaultLogName),
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
		AWSLoadBalancerControllerService: o.awsLoadBalancerControllerService(o.StormDB.From(constant.DefaultStormNodeAWSLoadBalanerController)),
		ArgoCD:                           o.argocdService(outputDir, o.StormDB.From(constant.DefaultStormNodeArgoCD), spin),
		ApplicationService:               o.applicationService(applicationsOutputDir, o.StormDB.From(constant.DefaultStormNodeApplications), spin),
		Certificate:                      o.certService(o.StormDB.From(constant.DefaultStormNodeCertificates)),
		Cluster:                          o.clusterService(outputDir, spin),
		Domain:                           o.domainService(o.StormDB.From(constant.DefaultStormNodeDomains)),
		ExternalDNS:                      o.externalDNSService(o.StormDB.From(constant.DefaultStormNodeExternalDNS)),
		ExternalSecrets:                  o.externalSecretsService(o.StormDB.From(constant.DefaultStormNodeExternalSecrets)),
		Github:                           o.githubService(ghClient, spin),
		Manifest:                         o.manifestService(o.StormDB.From(constant.DefaultStormNodeKubernetesManifest)),
		NameserverHandler:                o.nameserverHandlerService(ghClient, outputDir, spin),
		Parameter:                        o.paramService(o.StormDB.From(constant.DefaultStormNodeParameter)),
		Vpc:                              o.vpcService(o.StormDB.From(constant.DefaultStormNodeVpc)),
		IdentityManager:                  o.identityManagerService(o.StormDB.From(constant.DefaultStormNodeIdentityManager)),
		Autoscaler:                       o.autoscalerService(o.StormDB.From(constant.DefaultstormNodeAutoscaler)),
		Blockstorage:                     o.blockstorageService(o.StormDB.From(constant.DefaultStormNodeBlockStorage)),
		Monitoring:                       o.monitoringService(outputDir, o.StormDB.From(constant.DefaultStormNodeMonitoring), spin),
		Component:                        o.componentService(outputDir, o.StormDB.From(constant.DefaultStormNodeComponent), spin),
		Helm:                             o.helmService(o.StormDB.From(constant.DefaultStormNodeHelm)),
	}, nil
}

func (o *Okctl) helmService(node stormpkg.Node) client.HelmService {
	return clientCore.NewHelmService(
		rest.NewHelmAPI(o.restClient),
		storm.NewHelmState(node),
	)
}

func (o *Okctl) componentService(outputDir string, node stormpkg.Node, spin spinner.Spinner) client.ComponentService {
	return clientCore.NewComponentService(spin,
		rest.NewComponentAPI(o.restClient),
		clientFilesystem.NewComponentStore(
			clientFilesystem.Paths{
				OutputFile:         constant.DefaultPostgresOutputFile,
				CloudFormationFile: constant.DefaultPostgresCloudFormationFile,
				BaseDir:            path.Join(outputDir, constant.DefaultComponentBaseDir, constant.DefaultPostgresBaseDir),
			},
			o.FileSystem,
		),
		stateSaver.NewComponentState(o.RepoStateWithEnv),
		console.NewComponentReport(o.Err, spin),
		o.manifestService(node),
		o.CloudProvider,
	)
}

func (o *Okctl) managedPolicyService(node stormpkg.Node) client.ManagedPolicyService {
	return clientCore.NewManagedPolicyService(
		rest.NewManagedPolicyAPI(o.restClient),
		storm.NewManagedPolicyState(node),
	)
}

func (o *Okctl) serviceAccountService(node stormpkg.Node) client.ServiceAccountService {
	return clientCore.NewServiceAccountService(
		rest.NewServiceAccountAPI(o.restClient),
		storm.NewServiceAccountState(node),
	)
}

func (o *Okctl) monitoringService(outputDir string, node stormpkg.Node, spin spinner.Spinner) client.MonitoringService {
	monitoringDir := path.Join(outputDir, constant.DefaultMonitoringBaseDir)

	return clientCore.NewMonitoringService(
		rest.NewMonitoringAPI(o.restClient),
		clientFilesystem.NewMonitoringStore(
			clientFilesystem.Paths{
				OutputFile:  constant.DefaultHelmOutputsFile,
				ReleaseFile: constant.DefaultHelmReleaseFile,
				ChartFile:   constant.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, constant.DefaultTempoBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  constant.DefaultHelmOutputsFile,
				ReleaseFile: constant.DefaultHelmReleaseFile,
				ChartFile:   constant.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, constant.DefaultPromtailBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  constant.DefaultHelmOutputsFile,
				ReleaseFile: constant.DefaultHelmReleaseFile,
				ChartFile:   constant.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, constant.DefaultLokiBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile:  constant.DefaultHelmOutputsFile,
				ReleaseFile: constant.DefaultHelmReleaseFile,
				ChartFile:   constant.DefaultHelmChartFile,
				BaseDir:     path.Join(monitoringDir, constant.DefaultKubePromStackBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: constant.DefaultKubePromStackOutputsFile,
				BaseDir:    path.Join(monitoringDir, constant.DefaultKubePromStackBaseDir),
			},
			o.FileSystem,
		),
		stateSaver.NewMonitoringState(o.RepoStateWithEnv),
		console.NewMonitoringReport(o.Err, spin),
		o.certService(node),
		o.identityManagerService(node),
		o.manifestService(node),
		o.paramService(node),
		o.serviceAccountService(node),
		o.managedPolicyService(node),
		o.CloudProvider,
	)
}

func (o *Okctl) identityManagerService(node stormpkg.Node) client.IdentityManagerService {
	return clientCore.NewIdentityManagerService(
		rest.NewIdentityManagerAPI(o.restClient),
		storm.NewIdentityManager(node),
		o.certService(node),
	)
}

func (o *Okctl) argocdService(outputDir string, node stormpkg.Node, spin spinner.Spinner) client.ArgoCDService {
	argoBaseDir := path.Join(outputDir, constant.DefaultArgoCDBaseDir)

	argoService := clientCore.NewArgoCDService(
		o.identityManagerService(node),
		o.certService(node),
		o.manifestService(node),
		o.paramService(node),
		rest.NewArgoCDAPI(o.restClient),
		clientFilesystem.NewArgoCDStore(
			clientFilesystem.Paths{
				OutputFile:  constant.DefaultHelmOutputsFile,
				ReleaseFile: constant.DefaultHelmReleaseFile,
				ChartFile:   constant.DefaultHelmChartFile,
				BaseDir:     path.Join(outputDir, constant.DefaultHelmBaseDir),
			},
			clientFilesystem.Paths{
				OutputFile: constant.DefaultArgoOutputsFile,
				BaseDir:    argoBaseDir,
			},
			o.FileSystem,
		),
		console.NewArgoCDReport(o.Err, spin),
		stateSaver.NewArgoCDState(o.RepoStateWithEnv),
	)

	return argoService
}

func (o *Okctl) paramService(node stormpkg.Node) client.ParameterService {
	return clientCore.NewParameterService(
		rest.NewParameterAPI(o.restClient),
		storm.NewParameterState(node),
	)
}

func (o *Okctl) manifestService(node stormpkg.Node) client.ManifestService {
	return clientCore.NewManifestService(
		rest.NewManifestAPI(o.restClient),
		storm.NewManifestState(node),
	)
}

func (o *Okctl) nameserverHandlerService(ghClient githubClient.Githuber, _ string, spin spinner.Spinner) client.NameserverRecordDelegationService {
	return clientCore.NewNameserverHandlerService(
		ghClient,
		spin,
	)
}

func (o *Okctl) certService(node stormpkg.Node) client.CertificateService {
	return clientCore.NewCertificateService(
		rest.NewCertificateAPI(o.restClient),
		storm.NewCertificateState(node),
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

func (o *Okctl) vpcService(node stormpkg.Node) client.VPCService {
	return clientCore.NewVPCService(
		rest.NewVPCAPI(o.restClient),
		storm.NewVpcState(node),
	)
}

func (o *Okctl) clusterService(outputDir string, spin spinner.Spinner) client.ClusterService {
	return clientCore.NewClusterService(
		spin,
		rest.NewClusterAPI(o.restClient),
		clientFilesystem.NewClusterStore(
			clientFilesystem.Paths{
				ConfigFile: constant.DefaultClusterConfig,
				BaseDir:    path.Join(outputDir, constant.DefaultClusterBaseDir),
			},
			o.FileSystem,
		),
		console.NewClusterReport(o.Err, spin),
		stateSaver.NewClusterState(o.RepoStateWithEnv),
		o.CloudProvider,
		o.CredentialsProvider.Aws(),
	)
}

func (o *Okctl) autoscalerService(node stormpkg.Node) client.AutoscalerService {
	return clientCore.NewAutoscalerService(
		o.managedPolicyService(node),
		o.serviceAccountService(node),
		o.helmService(node),
	)
}

func (o *Okctl) blockstorageService(node stormpkg.Node) client.BlockstorageService {
	return clientCore.NewBlockstorageService(
		o.managedPolicyService(node),
		o.serviceAccountService(node),
		o.helmService(node),
		o.manifestService(node),
	)
}

func (o *Okctl) externalSecretsService(node stormpkg.Node) client.ExternalSecretsService {
	return clientCore.NewExternalSecretsService(
		o.managedPolicyService(node),
		o.serviceAccountService(node),
		o.helmService(node),
	)
}

func (o *Okctl) awsLoadBalancerControllerService(node stormpkg.Node) client.AWSLoadBalancerControllerService {
	return clientCore.NewAWSLoadBalancerControllerService(
		o.managedPolicyService(node),
		o.serviceAccountService(node),
		o.helmService(node),
	)
}

func (o *Okctl) applicationService(applicationOutputDir string, node stormpkg.Node, spin spinner.Spinner) client.ApplicationService {
	return clientCore.NewApplicationService(
		o.FileSystem,
		clientFilesystem.Paths{
			BaseDir: applicationOutputDir,
		},
		o.certService(node),
		clientFilesystem.NewApplicationStore(
			clientFilesystem.Paths{
				BaseDir: applicationOutputDir,
			},
			o.FileSystem,
		),
		console.NewApplicationReport(o.Out, spin),
	)
}

func (o *Okctl) domainService(node stormpkg.Node) client.DomainService {
	return clientCore.NewDomainService(
		rest.NewDomainAPI(o.restClient),
		storm.NewDomainState(node),
	)
}

func (o *Okctl) externalDNSService(node stormpkg.Node) client.ExternalDNSService {
	return clientCore.NewExternalDNSService(
		rest.NewExternalDNSAPI(o.restClient),
		storm.NewExternalDNSState(node),
		o.managedPolicyService(node),
		o.serviceAccountService(node),
	)
}

// KubeConfigStore returns an initialised kube config store
func (o *Okctl) KubeConfigStore() (api.KubeConfigStore, error) {
	appDir, err := o.GetUserDataDir()
	if err != nil {
		return nil, err
	}

	outputDir, err := o.GetRepoOutputDir(o.activeEnv)
	if err != nil {
		return nil, err
	}

	return filesystem.NewKubeConfigStore(
		o.CloudProvider,
		constant.DefaultClusterKubeConfig,
		path.Join(appDir, constant.DefaultCredentialsDirName, o.RepoStateWithEnv.GetClusterName()),
		constant.DefaultClusterConfig,
		path.Join(outputDir, constant.DefaultClusterBaseDir),
		o.FileSystem,
	), nil
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

	err = o.initialiseStorm()
	if err != nil {
		return err
	}

	o.restClient = rest.New(o.Debug, o.Err, o.ServerURL)

	homeDir, err := o.GetHomeDir()
	if err != nil {
		return err
	}

	appDir, err := o.GetUserDataDir()
	if err != nil {
		return err
	}

	clusterName := o.RepoStateWithEnv.GetClusterName()

	kubeConfigStore, err := o.KubeConfigStore()
	if err != nil {
		return err
	}

	o.kubeConfigStore = kubeConfigStore

	vpcService := core.NewVpcService(
		awsProvider.NewVpcCloud(o.CloudProvider),
	)

	clusterService := core.NewClusterService(
		run.NewClusterRun(
			o.Debug,
			kubeConfigStore,
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsConfig),
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsCredentials),
			o.BinariesProvider,
			o.CloudProvider,
		),
	)

	managedPolicyService := core.NewManagedPolicyService(awsProvider.NewManagedPolicyCloudProvider(o.CloudProvider))

	serviceAccountService := core.NewServiceAccountService(
		run.NewServiceAccountRun(
			o.Debug,
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsConfig),
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsCredentials),
			o.BinariesProvider,
		),
	)

	kubeService := core.NewKubeService(
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
			HelmPluginsDirectory: path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmPluginsDirectory),
			HelmRegistryConfig:   path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmRegistryConfig),
			HelmRepositoryConfig: path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryConfig),
			HelmRepositoryCache:  path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryCache),
			HelmBaseDir:          path.Join(appDir, constant.DefaultHelmBaseDir),
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
	)

	domainService := core.NewDomainService(
		awsProvider.NewDomainCloudProvider(o.CloudProvider),
	)

	certificateService := core.NewCertificateService(
		awsProvider.NewCertificateCloudProvider(o.CloudProvider),
	)

	parameterService := core.NewParameterService(
		awsProvider.NewParameterCloudProvider(o.CloudProvider),
	)

	componentService := core.NewComponentService(
		awsProvider.NewComponentCloudProvider(o.CloudProvider),
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
		Cluster:          clusterService,
		Vpc:              vpcService,
		ManagedPolicy:    managedPolicyService,
		ServiceAccount:   serviceAccountService,
		Helm:             helmService,
		Kube:             kubeService,
		Domain:           domainService,
		Certificate:      certificateService,
		Parameter:        parameterService,
		IdentityManager:  identityManagerService,
		ComponentService: componentService,
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

func (o *Okctl) initialiseStorm() error {
	outputDir, err := o.GetRepoOutputDir(o.activeEnv)
	if err != nil {
		return err
	}

	s, err := stormpkg.Open(path.Join(outputDir, constant.DefaultStormDBName), stormpkg.Codec(json.Codec))
	if err != nil {
		return err
	}

	o.StormDB = s

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
		return errors.E(err, "creating binaries fetcher", errors.Internal)
	}

	out := ioutil.Discard
	if o.Debug {
		out = o.Err
	}

	o.BinariesProvider = binaries.New(o.Logger, out, o.CredentialsProvider.Aws(), fetcher)

	return nil
}

// NewCloudProviderWithRegion create a cloud provider with a specific region
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
		constant.DefaultClusterAwsConfig,
		constant.DefaultClusterAwsCredentials,
		path.Join(appDir, constant.DefaultCredentialsDirName, o.RepoStateWithEnv.GetClusterName()),
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
		return fmt.Errorf("acquiring AWS authenticator: %w", err)
	}

	githubAuthenticator, err := o.getGithubAuthenticator()
	if err != nil {
		return fmt.Errorf("acquiring Github authenticator: %w", err)
	}

	o.CredentialsProvider = credentials.New(awsAuthenticator, githubAuthenticator)

	return nil
}
