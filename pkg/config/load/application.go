package load

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ErrOnAppDataNotFound returns a handler that errors
// when the app configuration cannot be loaded
func ErrOnAppDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		f, _ := c.GetAppDataPath()

		return &DataNotFoundErr{
			err: fmt.Errorf("failed to load app configuration from: %s", f),
		}
	}
}

// CreateOnAppDataNotFound will start an interactive survey
// that allows the end user to configure okctl for their
// use
// nolint
func CreateOnAppDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		home, err := c.GetHomeDir()
		if err != nil {
			return err
		}

		{
			doContinue := false
			prompt := &survey.Confirm{
				Message: "First time using okctl? Start guided configuration?",
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

		appDataPath, err := c.GetAppDataPath()
		if err != nil {
			return err
		}
		{
			doContinue := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Configuration will be written to: %s. Continue?", appDataPath),
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

		store := storage.NewFileSystemStorage(home)

		writer, err := store.Create(config.DefaultDir, config.DefaultConfig, 0644)
		if err != nil {
			return err
		}

		defer func() {
			err = writer.Close()
		}()

		data := application.New()

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
			"configuration_file": appDataPath,
		}).Info("cli configuration completed")

		c.AppData = data

		return nil
	}
}

func AppDataFromFlagsThenEnvVarsThenConfigFile(cmd *cobra.Command, notFoundFn DataNotFoundFn) config.DataLoaderFn {
	return buildAppDataLoader(notFoundFn, func(v *viper.Viper) {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			viperName := strings.ReplaceAll(f.Name, "-", "_")
			if !f.Changed && v.IsSet(viperName) {
				val := v.Get(viperName)
				_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			}
		})
	})
}

func buildAppDataLoader(notFoundFn DataNotFoundFn, viperCfg func(v *viper.Viper)) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		v := viper.New()
		v.SetFs(cfg.FileSystem)

		b, err := yaml.Marshal(application.New())
		if err != nil {
			return err
		}

		defaultAppData := bytes.NewReader(b)

		v.SetConfigType(config.DefaultConfigType)

		err = v.MergeConfig(defaultAppData)
		if err != nil {
			return err
		}

		configPath, err := cfg.GetAppDataDir()
		if err != nil {
			return err
		}

		v.AddConfigPath(configPath)
		v.SetConfigName(config.DefaultConfigName)

		err = v.MergeInConfig()
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

		v.AutomaticEnv()
		v.SetEnvPrefix(config.EnvPrefix)
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if viperCfg != nil {
			viperCfg(v)
		}

		appData := &application.Data{}

		err = v.Unmarshal(appData)
		if err != nil {
			return err
		}

		cfg.AppData = appData

		return nil
	}
}
