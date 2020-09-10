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

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/oslokommune/okctl/pkg/api/core/store/noop"

	"github.com/oslokommune/okctl/pkg/credentials/github"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api/core"
	awsProvider "github.com/oslokommune/okctl/pkg/api/core/cloudprovider/aws"
	"github.com/oslokommune/okctl/pkg/api/core/run"
	"github.com/oslokommune/okctl/pkg/api/core/store/filesystem"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
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

// Okctl stores all state required for invoking commands
type Okctl struct {
	*config.Config

	CloudProvider       v1alpha1.CloudProvider
	BinariesProvider    binaries.Provider
	CredentialsProvider credentials.Provider
}

// Initialise okctl for receiving requests
// nolint: funlen
func (o *Okctl) Initialise(env, awsAccountID string) error {
	err := o.EnableFileLog()
	if err != nil {
		return err
	}

	err = o.initialiseProviders(env, awsAccountID)
	if err != nil {
		return err
	}

	homeDir, err := o.GetHomeDir()
	if err != nil {
		return err
	}

	appDir, err := o.GetUserDataDir()
	if err != nil {
		return err
	}

	kubeConfigStore := filesystem.NewKubeConfigStore(
		config.DefaultClusterKubeConfig,
		path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env)),
		o.FileSystem,
	)

	vpcStore := noop.NewVpcStore()
	clusterStore := noop.NewClusterStore()
	managedPolicyStore := noop.NewManagedPolicyStore()
	serviceAccountStore := noop.NewServiceAccountStore()
	helmStore := noop.NewHelmStore()
	kubeStore := noop.NewKubeStore()
	domainStore := noop.NewDomainStore()
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
			path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env), config.DefaultClusterAwsConfig),
			path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env), config.DefaultClusterAwsCredentials),
			o.BinariesProvider,
			o.CloudProvider,
		),
	)

	managedPolicyService := core.NewManagedPolicyService(
		awsProvider.NewManagedPolicyCloudProvider(o.CloudProvider),
		managedPolicyStore,
	)

	serviceAccountService := core.NewServiceAccountService(
		serviceAccountStore,
		run.NewServiceAccountRun(
			o.Debug,
			path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env), config.DefaultClusterAwsConfig),
			path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env), config.DefaultClusterAwsCredentials),
			o.BinariesProvider,
		),
	)

	kubeService := core.NewKubeService(
		kubeStore,
		run.NewKubeRun(kubeConfigStore),
	)

	awsIamAuth, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return err
	}

	helmRun := run.NewHelmRun(
		helm.New(&helm.Config{
			HomeDir:              homeDir,
			Path:                 fmt.Sprintf("/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:%s", awsIamAuth.BinaryPath),
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
		domainStore,
	)

	certificateService := core.NewCertificateService(
		awsProvider.NewCertificateCloudProvider(o.CloudProvider),
		certificateStore,
	)

	parameterService := core.NewParameterService(
		awsProvider.NewParameterCloudProvider(o.CloudProvider),
		parameterStore,
	)

	services := core.Services{
		Cluster:        clusterService,
		Vpc:            vpcService,
		ManagedPolicy:  managedPolicyService,
		ServiceAccount: serviceAccountService,
		Helm:           helmService,
		Kube:           kubeService,
		Domain:         domainService,
		Certificate:    certificateService,
		Parameter:      parameterService,
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

// New returns a new okctl instance
func New() *Okctl {
	return &Okctl{
		Config: config.New(),
	}
}

// Binaries returns the application binaries
func (o *Okctl) Binaries() []state.Binary {
	return o.UserData.Binaries
}

// Host returns the host information
func (o *Okctl) Host() state.Host {
	return o.UserData.Host
}

// Username returns the username of the active user
func (o *Okctl) Username() string {
	return o.UserData.User.Username
}

// Region returns the AWS region of the repository
func (o *Okctl) Region() string {
	return o.RepoData.Region
}

// initialiseProviders knows how to create all required providers
func (o *Okctl) initialiseProviders(env, awsAccountID string) error {
	err := o.newCredentialsProvider(env, awsAccountID)
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
		true,
		o.Host(),
		o.Binaries(),
		storage.NewFileSystemStorage(userDataDir),
	)
	if err != nil {
		return errors.E(err, "failed to create binaries fetcher", errors.Internal)
	}

	o.BinariesProvider = binaries.New(ioutil.Discard, o.CredentialsProvider.Aws(), fetcher)

	return nil
}

// newCloudProvider creates a provider for running cloud operations
func (o *Okctl) newCloudProvider() error {
	c, err := cloud.New(o.Region(), o.CredentialsProvider.Aws())
	if err != nil {
		return err
	}

	o.CloudProvider = c.Provider

	return nil
}

// newCredentialsProvider knows how to load credentials
func (o *Okctl) newCredentialsProvider(env, awsAccountID string) error {
	if o.NoInput {
		return errors.E(errors.Errorf("we only support retrieving credentials interactively for now"), errors.Invalid)
	}

	appDir, err := o.GetUserDataDir()
	if err != nil {
		return err
	}

	authStore := aws.NewIniPersister(aws.NewFileSystemIniStorer(
		config.DefaultClusterAwsConfig,
		config.DefaultClusterAwsCredentials,
		path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env)),
		o.FileSystem,
	))

	defaultRing, err := keyring.DefaultKeyring()
	if err != nil {
		return err
	}

	k, err := keyring.New(defaultRing)
	if err != nil {
		return err
	}

	storedPassword, _ := k.Fetch(keyring.KeyTypeUserPassword)
	fn := func(username, password string) {
		// We do not handle this error, since we do not want the process to stop even if we cannot
		// save password
		_ = k.Store(keyring.KeyTypeUserPassword, password)
	}

	saml := aws.NewAuthSAML(awsAccountID, o.Region(), scrape.New(), aws.DefaultStsProvider, aws.Interactive(o.Username(), storedPassword, fn))

	gh := github.New(
		github.NewKeyringPersister(k),
		&http.Client{
			Timeout: 5 * time.Second, // nolint: gomnd
		},
		github.NewAuthDeviceFlow(github.DefaultGithubOauthClientID, github.RequiredScopes()),
	)

	o.CredentialsProvider = credentials.New(aws.New(authStore, saml), gh)

	return nil
}
