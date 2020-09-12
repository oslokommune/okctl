package load

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ErrOnRepoDataNotFound will simply error if no
// application data is found
func ErrOnRepoDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		f, _ := c.GetRepoDataPath()

		return &DataNotFoundErr{
			err: fmt.Errorf("failed to load repo configuration from: %s", f),
		}
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

		repoDataPath, err := c.GetRepoDataPath()
		if err != nil {
			return err
		}

		err = PromptContinue(fmt.Sprintf("GithubRepository configuration will be written to: %s. Continue?", repoDataPath), "user aborted configuration")
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

		err = c.WriteCurrentRepoData()
		if err != nil {
			return errors.Wrap(err, "failed to write current repo data")
		}

		c.Logger.WithFields(logrus.Fields{
			"configuration_file": repoDataPath,
		}).Info("repository configuration completed")

		return nil
	}
}

// RepoDataFromConfigFile defines the default behavior for loading configuration data
func RepoDataFromConfigFile(_ *cobra.Command, notFoundFn DataNotFoundFn) config.DataLoaderFn {
	return buildRepoDataLoader(notFoundFn, nil)
}

func buildRepoDataLoader(notFoundFn DataNotFoundFn, _ func(v *viper.Viper)) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		baseDir, err := cfg.GetRepoDir()
		if err != nil {
			return err
		}

		cfg.RepoState = &state.Repository{}

		exists, err := cfg.FileSystem.Exists(path.Join(baseDir, config.DefaultRepositoryConfig))
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
			GetStruct(config.DefaultRepositoryConfig, cfg.RepoState, store.FromYAML()).
			Do()
		if err != nil {
			return err
		}

		return nil
	}
}
