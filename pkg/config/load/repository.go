package load

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ErrOnRepoDataNotFound will simply error if no
// application data is found
func ErrOnRepoDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		f, _ := c.GetRepoStatePath()

		return &DataNotFoundErr{
			err: fmt.Errorf("failed to load repo configuration from: %s", f),
		}
	}
}

const repoMsg = `
We will now setup this repository for use with okctl by
creating a configuration file in it and for writing state
to, such as cloud formation templates, output from actions,
and more.

We will ask some questions about where and how this state
will be stored, together with your AWS account details.

`

// CreateOnRepoDataNotFoundWithNoUserInput configures okctl when a repository configuration was not found
func CreateOnRepoDataNotFoundWithNoUserInput(declaration *v1alpha1.Cluster) DataNotFoundFn {
	return func(c *config.Config) error {
		var err error

		repoDataPath, err := c.GetRepoStatePath()
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(c.Err, repoMsg)
		if err != nil {
			return err
		}

		data := state.NewRepository()

		data.Metadata = state.Metadata{
			Name:      declaration.Metadata.Name,
			Region:    declaration.Metadata.Region,
			OutputDir: declaration.Github.OutputPath,
		}

		c.RepoState = data

		repoDir, err := c.GetRepoDir()
		if err != nil {
			return err
		}

		_, err = store.NewFileSystem(repoDir, c.FileSystem).
			StoreStruct(constant.DefaultRepositoryStateFile, c.RepoState, store.ToYAML()).
			Do()
		if err != nil {
			return err
		}

		c.Logger.WithFields(logrus.Fields{
			"configuration_file": repoDataPath,
		}).Info("repository configuration completed")

		return nil
	}
}

// CreateOnRepoDataNotFound will start an interactive survey
// that allows the end user to configure okctl when a repository
// configuration was not found
func CreateOnRepoDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		var err error

		err = PromptContinue("First time using okctl with this repository? Start guided configuration?", "user aborted configuration")
		if err != nil {
			return err
		}

		repoDataPath, err := c.GetRepoStatePath()
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(c.Err, repoMsg)
		if err != nil {
			return err
		}

		err = PromptContinue(fmt.Sprintf("Repository configuration will be written to: %s. Continue?", repoDataPath), "user aborted configuration")
		if err != nil {
			return err
		}

		data := state.NewRepository()

		cfg, err := ask.New().RepositoryConfig()
		if err != nil {
			return err
		}

		data.Metadata = state.Metadata{
			Name:      cfg.Name,
			Region:    cfg.Region,
			OutputDir: cfg.BaseDir,
		}

		c.RepoState = data

		repoDir, err := c.GetRepoDir()
		if err != nil {
			return err
		}

		_, err = store.NewFileSystem(repoDir, c.FileSystem).
			StoreStruct(constant.DefaultRepositoryStateFile, c.RepoState, store.ToYAML()).
			Do()
		if err != nil {
			return err
		}

		c.Logger.WithFields(logrus.Fields{
			"configuration_file": repoDataPath,
		}).Info("repository configuration completed")

		return nil
	}
}

// RepoDataFromConfigFile defines the default behavior for loading configuration data
func RepoDataFromConfigFile(notFoundFn DataNotFoundFn) config.DataLoaderFn {
	// Here we should be loading all the definitions in the `.okctl/` directory together
	// with the .okctl.state file. Once these are loaded we can validate the `.okctl/*.yml` files
	// and start working on matching the declared state with the actual state.
	//
	// We probably need to get credentials for all the different AWS thingies before we get started?
	// Perhaps the apply thingy should specify the cluster? So we don't operate on all declared states?
	return buildRepoDataLoader(notFoundFn, nil)
}

func buildRepoDataLoader(notFoundFn DataNotFoundFn, _ func(v *viper.Viper)) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		baseDir, err := cfg.GetRepoDir()
		if err != nil {
			return err
		}

		cfg.RepoState = &state.Repository{}

		exists, err := cfg.FileSystem.Exists(path.Join(baseDir, constant.DefaultRepositoryStateFile))
		if err != nil {
			return err
		}

		if !exists {
			err = notFoundFn(cfg)
			if err != nil {
				return err
			}
		}

		_, err = store.NewFileSystem(baseDir, cfg.FileSystem).
			GetStruct(constant.DefaultRepositoryStateFile, cfg.RepoState, store.FromYAML()).
			Do()
		if err != nil {
			return err
		}

		return nil
	}
}
