package okctl

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/core"
	cld "github.com/oslokommune/okctl/pkg/api/core/cloud"
	"github.com/oslokommune/okctl/pkg/api/core/exe"
	"github.com/oslokommune/okctl/pkg/api/core/store"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/oslokommune/okctl/pkg/credentials/login"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/storage/state"
)

type ServiceProvider interface {
	ClusterService() api.ClusterService
}

// Okctl stores all state required for invoking commands
type Okctl struct {
	*config.Config

	CloudProvider       v1alpha1.CloudProvider
	BinariesProvider    binaries.Provider
	CredentialsProvider credentials.Provider
	PersisterProvider   state.PersisterProvider
}

func (o *Okctl) Initialise(env, awsAccountID string) error {
	err := o.initialiseProviders(env, awsAccountID)
	if err != nil {
		return err
	}

	clusterService := core.NewClusterService(
		store.NewClusterStore(o.PersisterProvider),
		cld.NewCluster(o.CloudProvider),
		exe.NewClusterExe(o.BinariesProvider),
	)

	services := core.Services{
		Cluster: clusterService,
	}

	endpoints := core.GenerateEndpoints(services, core.InstrumentEndpoints(o.Logger))

	handlers := core.MakeHandlers(o.Format(), endpoints)

	router := http.NewServeMux()
	router.Handle("/", core.AttachRoutes(handlers))

	server := &http.Server{
		Handler: router,
		Addr:    o.Destination,
	}

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

func (o *Okctl) Binaries() []application.Binary {
	return o.AppData.Binaries
}

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
	err := o.newCredentialsProvider(awsAccountID)
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

	err = o.newPersisterProvider(env)
	if err != nil {
		return err
	}

	return nil
}

// newPersisterProvider creates a provider for persisting state
func (o *Okctl) newPersisterProvider(env string) error {
	appDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	appOpts := state.AppStoreOpts{
		Opts: state.Opts{
			BaseDir:    appDir,
			ConfigFile: config.DefaultConfig,
			Defaults:   map[string]string{},
		},
		State: nil,
	}

	repoDir, err := o.GetRepoDir()
	if err != nil {
		return err
	}

	outputDir, err := o.GetRepoOutputDir(env)
	if err != nil {
		return err
	}

	repoOpts := state.RepoStoreOpts{
		Opts: state.Opts{
			BaseDir:    repoDir,
			ConfigFile: config.DefaultRepositoryConfig,
			Defaults: map[string]string{
				"cluster_config": path.Join(outputDir, config.DefaultClusterBaseDir, config.DefaultClusterConfig),
			},
		},
	}

	o.PersisterProvider = state.New(repoOpts, appOpts)

	return nil
}

// newBinariesProvider creates a provider for loading binaries
func (o *Okctl) newBinariesProvider() error {
	appDataDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	str := storage.NewFileSystemStorage(appDataDir)

	stagers, err := fetch.New(o.Host(), str).FromConfig(true, o.Binaries())
	if err != nil {
		return err
	}

	bin := binaries.New(o.Out, o.CredentialsProvider, stagers)

	o.BinariesProvider = bin

	return nil
}

// newCloudProvider creates a provider for running cloud operations
func (o *Okctl) newCloudProvider() error {
	c, err := cloud.New(o.Region(), o.CredentialsProvider)
	if err != nil {
		return err
	}

	o.CloudProvider = c.Provider

	return nil
}

// newCredentialsProvider knows how to load credentials
func (o *Okctl) newCredentialsProvider(awsAccountID string) error {
	if o.NoInput {
		return fmt.Errorf("we only support retrieving credentials interactively for now")
	}

	l, err := login.Interactive(awsAccountID, o.Region(), o.Username())
	if err != nil {
		return err
	}

	o.CredentialsProvider = credentials.New(l)

	return nil
}
