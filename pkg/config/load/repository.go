package load

import (
	"bytes"
	"fmt"
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

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
// nolint
func CreateOnRepoDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		repo, err := c.GetRepoDir()
		if err != nil {
			return err
		}

		{
			doContinue := false
			prompt := &survey.Confirm{
				Message: "First time using okctl with this repository? Start guided configuration?",
				Default: true,
			}

			err = survey.AskOne(prompt, &doContinue)
			if err != nil {
				return err
			}

			if !doContinue {
				return fmt.Errorf("no configuration file found, and user did not want to start a configuration")
			}
		}

		repoDataPath, err := c.GetRepoDataPath()
		if err != nil {
			return err
		}
		{
			doContinue := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Repository configuration will be written to: %s. Continue?", repoDataPath),
				Default: true,
			}

			err = survey.AskOne(prompt, &doContinue)
			if err != nil {
				return err
			}

			if !doContinue {
				return fmt.Errorf("no configuration file found, and user did not accept configuration file location")
			}
		}

		store := storage.NewFileSystemStorage(repo)

		writer, err := store.Create("", config.DefaultRepositoryConfig, 0644)
		if err != nil {
			return err
		}

		defer func() {
			err = writer.Close()
		}()

		data := repository.New()

		err = data.Survey()
		if err != nil {
			return err
		}

		b, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, bytes.NewReader(b))
		if err != nil {
			return err
		}

		c.Logger.WithFields(logrus.Fields{
			"configuration_file": repoDataPath,
		}).Info("repository configuration completed")

		c.RepoData = data

		return nil
	}
}

func RepoDataFromConfigFile(_ *cobra.Command, notFoundFn DataNotFoundFn) config.DataLoaderFn {
	return buildRepoDataLoader(notFoundFn, nil)
}

func buildRepoDataLoader(notFoundFn DataNotFoundFn, viperCfg func(v *viper.Viper)) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		baseDir, err := cfg.GetRepoDir()
		if err != nil {
			return err
		}

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

		repoData := &repository.Data{}

		err = v.Unmarshal(repoData)
		if err != nil {
			return err
		}

		cfg.RepoData = repoData

		return nil
	}
}
