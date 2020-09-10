package load

import (
	"fmt"

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

		data.Name = cfg.Name
		data.Region = cfg.Region
		data.OutputDir = cfg.BaseDir

		c.RepoData = data

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

func buildRepoDataLoader(notFoundFn DataNotFoundFn, viperCfg func(v *viper.Viper)) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		baseDir, err := cfg.GetRepoDir()
		if err != nil {
			return err
		}

		cfg.RepoData = &state.Repository{}

		v := viper.New()
		v.SetFs(cfg.FileSystem.Fs)
		v.SetConfigName(config.DefaultRepositoryConfigName)
		v.SetConfigType(config.DefaultRepositoryConfigType)
		v.AddConfigPath(baseDir)

		err = v.ReadInConfig()
		if err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				err = notFoundFn(cfg)
				if err != nil {
					return err
				}
			default:
				return err
			}
		}

		if viperCfg != nil {
			viperCfg(v)
		}

		err = v.Unmarshal(cfg.RepoData)
		if err != nil {
			return err
		}

		return nil
	}
}
