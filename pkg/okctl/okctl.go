// Package okctl implements the core logic for creating providers
// and loading configuration state
package okctl

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api/core"
	cld "github.com/oslokommune/okctl/pkg/api/core/cloud"
	"github.com/oslokommune/okctl/pkg/api/core/exe"
	"github.com/oslokommune/okctl/pkg/api/core/store/filesystem"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/credentials/aws/scrape"
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

	appDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	repoDir, err := o.GetRepoDir()
	if err != nil {
		return err
	}

	outputDir, err := o.GetRepoOutputDir(env)
	if err != nil {
		return err
	}

	vpcStore := filesystem.NewVpcStore(
		config.DefaultVpcOutputs,
		config.DefaultVpcCloudFormationTemplate,
		path.Join(outputDir, config.DefaultVpcBaseDir),
		o.FileSystem,
	)

	kubeConfigStore := filesystem.NewKubeConfigStore(
		config.DefaultClusterKubeConfig,
		path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env)),
		o.FileSystem,
	)

	clusterConfigStore := filesystem.NewClusterConfigStore(
		config.DefaultClusterConfig,
		path.Join(outputDir, config.DefaultClusterBaseDir),
		o.FileSystem,
	)

	clusterStore := filesystem.NewClusterStore(
		config.DefaultRepositoryConfig,
		repoDir,
		o.FileSystem,
		o.RepoData,
	)

	vpcService := core.NewVpcService(
		cld.NewVpcCloud(o.CloudProvider),
		vpcStore,
	)

	clusterConfigService := core.NewClusterConfigService(
		clusterConfigStore,
		vpcStore,
	)

	clusterService := core.NewClusterService(
		clusterStore,
		clusterConfigStore,
		kubeConfigStore,
		exe.NewClusterExe(
			o.Debug,
			path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env), config.DefaultClusterAwsConfig),
			path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env), config.DefaultClusterAwsCredentials),
			o.BinariesProvider,
		),
	)

	services := core.Services{
		Cluster:       clusterService,
		ClusterConfig: clusterConfigService,
		Vpc:           vpcService,
	}

	endpoints := core.GenerateEndpoints(services, core.InstrumentEndpoints(o.Logger))

	handlers := core.MakeHandlers(o.Format(), endpoints)

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
func (o *Okctl) Binaries() []application.Binary {
	return o.AppData.Binaries
}

// Host returns the host information
func (o *Okctl) Host() application.Host {
	return o.AppData.Host
}

// Username returns the username of the active user
func (o *Okctl) Username() string {
	return o.AppData.User.Username
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
	appDataDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	fetcher, err := fetch.New(
		true,
		o.Host(),
		o.Binaries(),
		storage.NewFileSystemStorage(appDataDir),
	)
	if err != nil {
		return errors.E(err, "failed to create binaries fetcher", errors.Internal)
	}

	o.BinariesProvider = binaries.New(o.Out, o.CredentialsProvider.Aws(), fetcher)

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

	appDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	authStore := aws.NewIniPersister(aws.NewFileSystemIniStorer(
		config.DefaultClusterAwsConfig,
		config.DefaultClusterAwsCredentials,
		path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(env)),
		o.FileSystem,
	))

	saml := aws.NewAuthSAML(awsAccountID, o.Region(), scrape.New(), aws.DefaultStsProvider, aws.Interactive(o.Username()))

	o.CredentialsProvider = credentials.New(aws.New(authStore, saml))

	return nil
}
