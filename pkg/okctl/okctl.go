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

	"github.com/oslokommune/okctl/pkg/client/core/state/direct"

	"github.com/oslokommune/okctl/pkg/breeze"

	"github.com/oslokommune/okctl/pkg/client/core/state/storm"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/context"

	"github.com/logrusorgru/aurora/v3"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	clientDirectAPI "github.com/oslokommune/okctl/pkg/client/core/api/direct"
	"github.com/oslokommune/okctl/pkg/client/core/api/rest"
	githubClient "github.com/oslokommune/okctl/pkg/github"

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

	DB breeze.Client

	restClient      *rest.HTTPClient
	kubeConfigStore api.KubeConfigStore
}

// Initialise okctl
func (o *Okctl) Initialise() error {
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

%w`

	return func(err error) error {
		return fmt.Errorf(errMsg,
			command,
			path.Join(userDir, constant.DefaultLogDir, constant.DefaultLogName),
			aurora.Bold("#kjøremiljø-support"),
			err,
		)
	}
}

// StateNodes returns the initialised state nodes
func (o *Okctl) StateNodes() *clientCore.StateNodes {
	return &clientCore.StateNodes{
		ArgoCD:              o.DB.From(constant.DefaultStormNodeArgoCD),
		Certificate:         o.DB.From(constant.DefaultStormNodeCertificates),
		Cluster:             o.DB.From(constant.DefaultStormNodeCluster),
		Domain:              o.DB.From(constant.DefaultStormNodeDomains),
		ExternalDNS:         o.DB.From(constant.DefaultStormNodeExternalDNS),
		Github:              o.DB.From(constant.DefaultStormNodeGithub),
		Manifest:            o.DB.From(constant.DefaultStormNodeKubernetesManifest),
		Parameter:           o.DB.From(constant.DefaultStormNodeParameter),
		Vpc:                 o.DB.From(constant.DefaultStormNodeVpc),
		IdentityManager:     o.DB.From(constant.DefaultStormNodeIdentityManager),
		Monitoring:          o.DB.From(constant.DefaultStormNodeMonitoring),
		Component:           o.DB.From(constant.DefaultStormNodeComponent),
		Helm:                o.DB.From(constant.DefaultStormNodeHelm),
		ManagedPolicy:       o.DB.From(constant.DefaultStormNodeManagedPolicy),
		ServiceAccount:      o.DB.From(constant.DefaultStormNodeServiceAccount),
		ContainerRepository: o.DB.From(constant.DefaultStormNodeContainerRepository),
		Upgrade:             o.DB.From(constant.DefaultStormNodeUpgrade),
	}
}

// StateHandlers returns the initialised state handlers
func (o *Okctl) StateHandlers(nodes *clientCore.StateNodes) *clientCore.StateHandlers {
	helmClient := rest.NewHelmAPI(o.restClient)

	return &clientCore.StateHandlers{
		Helm:                      storm.NewHelmState(nodes.Helm),
		ManagedPolicy:             storm.NewManagedPolicyState(nodes.ManagedPolicy),
		ServiceAccount:            storm.NewServiceAccountState(nodes.ServiceAccount),
		Certificate:               storm.NewCertificateState(nodes.Certificate),
		IdentityManager:           storm.NewIdentityManager(nodes.IdentityManager),
		Github:                    storm.NewGithubState(nodes.Github),
		Manifest:                  storm.NewManifestState(nodes.Manifest),
		Vpc:                       storm.NewVpcState(nodes.Vpc),
		Parameter:                 storm.NewParameterState(nodes.Parameter),
		Domain:                    storm.NewDomainState(nodes.Domain),
		ExternalDNS:               storm.NewExternalDNSState(nodes.ExternalDNS),
		Cluster:                   storm.NewClusterState(nodes.Cluster),
		Component:                 storm.NewComponentState(nodes.Component),
		Monitoring:                storm.NewMonitoringState(nodes.Monitoring),
		ArgoCD:                    storm.NewArgoCDState(nodes.ArgoCD),
		ContainerRepository:       storm.NewContainerRepositoryState(nodes.ContainerRepository),
		Loki:                      direct.NewLokiState(o.Declaration.Metadata, helmClient),
		Promtail:                  direct.NewPromtailState(o.Declaration.Metadata, helmClient),
		Tempo:                     direct.NewTempoState(o.Declaration.Metadata, helmClient),
		Autoscaler:                direct.NewAutoscalerState(o.Declaration.Metadata, helmClient),
		AWSLoadBalancerController: direct.NewAWSLoadBalancerState(o.Declaration.Metadata, helmClient),
		Blockstorage:              direct.NewBlockstorageState(o.Declaration.Metadata, helmClient),
		ExternalSecrets:           direct.NewExternalSecretsState(o.Declaration.Metadata, helmClient),
		Upgrade:                   storm.NewUpgradesState(nodes.Upgrade),
	}
}

// ClientServices returns the initialised client-side services
// nolint: funlen
func (o *Okctl) ClientServices(handlers *clientCore.StateHandlers) (*clientCore.Services, error) {
	absoluteRepositoryPath, err := o.GetRepoDir()
	if err != nil {
		return nil, err
	}

	applicationsOutputDir, err := o.GetRepoApplicationsOutputDir()
	if err != nil {
		return nil, err
	}

	ghClient, err := githubClient.New(o.Ctx, o.CredentialsProvider.Github())
	if err != nil {
		return nil, err
	}

	helmService := clientCore.NewHelmService(
		rest.NewHelmAPI(o.restClient),
		handlers.Helm,
	)

	managedPolicyService := clientCore.NewManagedPolicyService(
		rest.NewManagedPolicyAPI(o.restClient),
		handlers.ManagedPolicy,
	)

	serviceAccountService := clientCore.NewServiceAccountService(
		rest.NewServiceAccountAPI(o.restClient),
		handlers.ServiceAccount,
	)

	certificateService := clientCore.NewCertificateService(
		rest.NewCertificateAPI(o.restClient),
		handlers.Certificate,
	)

	usProvider, err := o.getUsEastOneProvider()
	if err != nil {
		return nil, errors.New("Unable to get certificate cloud provider")
	}

	identityManagerService := clientCore.NewIdentityManagerService(
		clientDirectAPI.NewIdentityManagerAPI(core.NewIdentityManagerService(
			awsProvider.NewIdentityManagerCloudProvider(o.CloudProvider),
			awsProvider.NewCertificateCloudProvider(usProvider),
		)),
		handlers.IdentityManager,
		certificateService,
	)

	githubService := clientCore.NewGithubService(
		rest.NewGithubAPI(
			rest.NewParameterAPI(o.restClient),
			ghClient,
		),
		handlers.Github,
	)

	autoscalerService := clientCore.NewAutoscalerService(
		managedPolicyService,
		serviceAccountService,
		helmService,
	)

	manifestService := clientCore.NewManifestService(
		rest.NewManifestAPI(o.restClient),
		handlers.Manifest,
	)

	applicationManifestService := clientCore.NewApplicationManifestService(
		o.FileSystem,
		applicationsOutputDir,
	)

	blockstorageService := clientCore.NewBlockstorageService(
		managedPolicyService,
		serviceAccountService,
		helmService,
		manifestService,
	)

	vpcService := clientCore.NewVPCService(
		rest.NewVPCAPI(o.restClient),
		handlers.Vpc,
	)

	paramService := clientCore.NewParameterService(
		rest.NewParameterAPI(o.restClient),
		handlers.Parameter,
	)

	externalSecretsService := clientCore.NewExternalSecretsService(
		managedPolicyService,
		serviceAccountService,
		helmService,
	)

	domainAPIService := core.NewDomainService(
		awsProvider.NewDomainCloudProvider(o.CloudProvider),
	)

	domainService := clientCore.NewDomainService(
		clientDirectAPI.NewDomainAPI(domainAPIService),
		handlers.Domain,
	)

	externalDNSService := clientCore.NewExternalDNSService(
		clientDirectAPI.NewExternalDNSAPI(core.NewKubeService(
			run.NewKubeRun(o.CloudProvider, o.CredentialsProvider.Aws()),
		)),
		handlers.ExternalDNS,
		managedPolicyService,
		serviceAccountService,
	)

	awsLoadBalancerControllerService := clientCore.NewAWSLoadBalancerControllerService(
		managedPolicyService,
		serviceAccountService,
		helmService,
	)

	appDir, err := o.GetUserDataDir()
	if err != nil {
		return nil, err
	}

	kubeConfigStore, err := o.KubeConfigStore()
	if err != nil {
		return nil, err
	}

	clusterName := o.Declaration.Metadata.Name

	clusterServiceCore := core.NewClusterService(
		run.NewClusterRun(
			o.Debug,
			kubeConfigStore,
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsConfig),
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsCredentials),
			o.BinariesProvider,
			o.CloudProvider,
		),
		o.CloudProvider,
	)

	clusterService := clientCore.NewClusterService(
		clientDirectAPI.NewClusterAPI(clusterServiceCore),
		handlers.Cluster,
		o.CloudProvider,
		o.CredentialsProvider.Aws(),
	)

	componentService := clientCore.NewComponentService(
		clientDirectAPI.NewComponentAPI(core.NewComponentService(
			awsProvider.NewComponentCloudProvider(o.CloudProvider),
		)),
		handlers.Component,
		manifestService,
		o.CloudProvider,
	)

	applicationPostgresService := clientCore.NewApplicationPostgresService(
		applicationManifestService,
		componentService,
		rest.NewSecurityGroupAPI(o.restClient),
		vpcService,
		rest.NewApplicationPostgresIntegrationAPI(o.restClient),
		clusterService,
	)

	monitoringService := clientCore.NewMonitoringService(
		handlers.Monitoring,
		helmService,
		certificateService,
		identityManagerService,
		manifestService,
		paramService,
		serviceAccountService,
		managedPolicyService,
		o.CloudProvider,
	)

	argocdService := clientCore.NewArgoCDService(
		identityManagerService,
		certificateService,
		manifestService,
		paramService,
		helmService,
		handlers.ArgoCD,
	)

	applicationService := clientCore.NewApplicationService(
		o.FileSystem,
		certificateService,
		applicationManifestService,
		absoluteRepositoryPath,
	)

	nameserverService := clientCore.NewNameserverHandlerService(ghClient)

	containerRepositoryService := clientCore.NewContainerRepositoryService(
		clientDirectAPI.NewContainerRepositoryAPI(core.NewContainerRepositoryService(
			awsProvider.NewContainerRepositoryCloudProvider(o.CloudProvider),
		)),
		handlers.ContainerRepository,
		o.CloudProvider,
	)

	services := &clientCore.Services{
		AWSLoadBalancerControllerService: awsLoadBalancerControllerService,
		ArgoCD:                           argocdService,
		ApplicationService:               applicationService,
		ApplicationManifestService:       applicationManifestService,
		ApplicationPostgresService:       applicationPostgresService,
		Certificate:                      certificateService,
		Cluster:                          clusterService,
		Domain:                           domainService,
		ExternalDNS:                      externalDNSService,
		ExternalSecrets:                  externalSecretsService,
		Github:                           githubService,
		Manifest:                         manifestService,
		NameserverHandler:                nameserverService,
		Parameter:                        paramService,
		Vpc:                              vpcService,
		IdentityManager:                  identityManagerService,
		Autoscaler:                       autoscalerService,
		Blockstorage:                     blockstorageService,
		Monitoring:                       monitoringService,
		Component:                        componentService,
		Helm:                             helmService,
		ManagedPolicy:                    managedPolicyService,
		ServiceAccount:                   serviceAccountService,
		ContainerRepository:              containerRepositoryService,
	}

	return services, nil
}

// KubeConfigStore returns an initialised kube config store
func (o *Okctl) KubeConfigStore() (api.KubeConfigStore, error) {
	appDir, err := o.GetUserDataDir()
	if err != nil {
		return nil, err
	}

	return filesystem.NewKubeConfigStore(
		o.CloudProvider,
		constant.DefaultClusterKubeConfig,
		path.Join(appDir, constant.DefaultCredentialsDirName, o.Declaration.Metadata.Name),
		o.StateHandlers(o.StateNodes()).Cluster,
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

	err = o.initialiseBreeze()
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

	kubeConfigStore, err := o.KubeConfigStore()
	if err != nil {
		return err
	}

	o.kubeConfigStore = kubeConfigStore

	vpcService := core.NewVpcService(
		awsProvider.NewVpcCloud(o.CloudProvider),
	)

	clusterName := o.Declaration.Metadata.Name

	clusterService := core.NewClusterService(
		run.NewClusterRun(
			o.Debug,
			kubeConfigStore,
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsConfig),
			path.Join(appDir, constant.DefaultCredentialsDirName, clusterName, constant.DefaultClusterAwsCredentials),
			o.BinariesProvider,
			o.CloudProvider,
		),
		o.CloudProvider,
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

	containerRepositoryService := core.NewContainerRepositoryService(
		awsProvider.NewContainerRepositoryCloudProvider(o.CloudProvider),
	)

	usProvider, err := o.getUsEastOneProvider()
	if err != nil {
		return errors.New("Unable to get certificate cloud provider")
	}

	identityManagerService := core.NewIdentityManagerService(
		awsProvider.NewIdentityManagerCloudProvider(o.CloudProvider),
		awsProvider.NewCertificateCloudProvider(usProvider),
	)

	securityGroupService := core.NewSecurityGroupService(
		awsProvider.NewSecurityGroupCloudProvider(o.CloudProvider),
	)

	services := core.Services{
		Cluster:                    clusterService,
		Vpc:                        vpcService,
		ManagedPolicy:              managedPolicyService,
		ServiceAccount:             serviceAccountService,
		Helm:                       helmService,
		Kube:                       kubeService,
		Domain:                     domainService,
		Certificate:                certificateService,
		Parameter:                  parameterService,
		IdentityManager:            identityManagerService,
		ComponentService:           componentService,
		ContainerRepositoryService: containerRepositoryService,
		SecurityGroupService:       securityGroupService,
	}

	endpoints := core.GenerateEndpoints(services, core.InstrumentEndpoints(o.Logger))

	handlers := core.MakeHandlers(core.EncodeJSONResponse, endpoints)

	router := http.NewServeMux()
	router.Handle("/", core.AttachRoutes(handlers))
	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

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

	return o.waitForServer()
}

// waitForServer waits for the http.Server to become active
func (o *Okctl) waitForServer() error {
	const (
		serverHealthTimeoutInSec = 5
		serverHealthIntervalInMs = 100
	)

	timeout := time.After(serverHealthTimeoutInSec * time.Second)
	tick := time.NewTicker(serverHealthIntervalInMs * time.Millisecond)

	for {
		select {
		case <-timeout:
			return errors.New("timed out waiting for server")
		case <-tick.C:
			r, err := http.Get(fmt.Sprintf("%s%s", o.ServerBaseURL, "health"))
			if err != nil {
				return err
			}

			defer func() {
				_ = r.Body.Close()
			}()

			if r.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
}

func (o *Okctl) initialiseBreeze() error {
	outputDir, err := o.GetRepoOutputDir()
	if err != nil {
		return err
	}

	db := breeze.New(path.Join(outputDir, constant.DefaultStormDBName))

	o.DB = db

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
	c, err := cloud.New(o.Declaration.Metadata.Region, o.CredentialsProvider.Aws())
	if err != nil {
		return err
	}

	o.CloudProvider = c.Provider

	return nil
}

func (o *Okctl) getAWSAuthenticator() (*aws.Auth, error) {
	if o.AWSCredentialsType == context.AWSCredentialsTypeAccessKey {
		return aws.New(aws.NewInMemoryStorage(), aws.NewAuthEnvironment(o.Declaration.Metadata.Region, os.Getenv)), nil
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
		path.Join(appDir, constant.DefaultCredentialsDirName, o.Declaration.Metadata.Name),
		o.FileSystem,
	))

	return aws.New(authStore, aws.NewAuthSAML(
		o.Declaration.Metadata.AccountID,
		o.Declaration.Metadata.Region,
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

// When creating a certificate for a CloudFront distribution, we
// need to create the certificate in us-east-1
func (o *Okctl) getUsEastOneProvider() (v1alpha1.CloudProvider, error) {
	provider, err := o.NewCloudProviderWithRegion("us-east-1")
	if err != nil {
		return nil, err
	}

	return provider, nil
}
