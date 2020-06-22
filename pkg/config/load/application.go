package load

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
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
func CreateOnAppDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		var err error

		err = PromptContinue(
			"First time using okctl? Start guided configuration?",
			"user aborted interactive configuration",
		)
		if err != nil {
			return err
		}

		appDataPath, err := c.GetAppDataPath()
		if err != nil {
			return err
		}

		err = PromptContinue(
			fmt.Sprintf("Configuration will be written to: %s. Continue?", appDataPath),
			"user aborted interactive configuration",
		)
		if err != nil {
			return err
		}

		data, err := application.New().Survey()
		if err != nil {
			return err
		}

		b, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		err = c.WriteAppData(b)
		if err != nil {
			return err
		}

		c.AppData = data

		c.Logger.WithFields(logrus.Fields{
			"configuration_file": appDataPath,
		}).Info("cli configuration completed")

		return nil
	}
}

type LoaderFn func(cfg *config.Config, v *viper.Viper) error

func AppDataFromFlagsEnvConfigDefaults(cmd *cobra.Command, notFoundFn DataNotFoundFn) config.DataLoaderFn {
	return buildAppDataLoader(
		loadDefaultAppData,
		loadStoredAppData(notFoundFn),
		loadEnvAppData,
		loadFlagsAppData(cmd),
	)
}

func loadDefaultAppData(_ *config.Config, v *viper.Viper) error {
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

	return nil
}

func loadStoredAppData(notFoundFn DataNotFoundFn) LoaderFn {
	return func(cfg *config.Config, v *viper.Viper) error {
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

		return nil
	}
}

func loadEnvAppData(_ *config.Config, v *viper.Viper) error {
	v.AutomaticEnv()
	v.SetEnvPrefix(config.EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}

func loadFlagsAppData(cmd *cobra.Command) LoaderFn {
	return func(cfg *config.Config, v *viper.Viper) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			viperName := strings.ReplaceAll(f.Name, "-", "_")
			if !f.Changed && v.IsSet(viperName) {
				val := v.Get(viperName)
				_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			}
		})

		return nil
	}
}

func buildAppDataLoader(loaders ...LoaderFn) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		var err error

		v := viper.New()
		v.SetFs(cfg.FileSystem.Fs)

		for _, fn := range loaders {
			err = fn(cfg, v)
			if err != nil {
				return err
			}
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
