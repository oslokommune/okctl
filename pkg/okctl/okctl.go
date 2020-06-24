package okctl

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/oslokommune/okctl/pkg/credentials/login"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/storage/state"
)

// Okctl stores all state required for invoking commands
type Okctl struct {
	*config.Config

	CloudProvider       v1alpha1.CloudProvider
	BinariesProvider    binaries.Provider
	CredentialsProvider credentials.Provider
	PersisterProvider   state.PersisterProvider
}

// New returns a new okctl instance
func New() *Okctl {
	return &Okctl{
		Config: config.New(),
	}
}

// Username returns the username of the active user
func (o *Okctl) Username() string {
	return o.AppData.User.Username
}

// Region returns the AWS region of the repository
func (o *Okctl) Region() string {
	return o.RepoData.Region
}

// NewProviders knows how to create all required providers
func (o *Okctl) NewProviders(env, awsAccountID string) error {
	err := o.NewCredentialsProvider(awsAccountID)
	if err != nil {
		return err
	}

	err = o.NewBinariesProvider()
	if err != nil {
		return err
	}

	err = o.NewCloudProvider()
	if err != nil {
		return err
	}

	err = o.NewPersisterProvider(env)
	if err != nil {
		return err
	}

	return nil
}

// NewPersisterProvider creates a provider for persisting state
func (o *Okctl) NewPersisterProvider(env string) error {
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

// NewBinariesProvider creates a provider for loading binaries
func (o *Okctl) NewBinariesProvider() error {
	appDataDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(appDataDir)

	stagers, err := fetch.New(o.AppData.Host, store).FromConfig(true, o.AppData.Binaries)
	if err != nil {
		return err
	}

	bin := binaries.New(o.Out, o.CredentialsProvider, stagers)

	o.BinariesProvider = bin

	return nil
}

// NewCloudProvider creates a provider for running cloud operations
func (o *Okctl) NewCloudProvider() error {
	c, err := cloud.New(o.Region(), o.CredentialsProvider)
	if err != nil {
		return err
	}

	o.CloudProvider = c.Provider

	return nil
}

// NewCredentialsProvider knows how to load credentials
func (o *Okctl) NewCredentialsProvider(awsAccountID string) error {
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
