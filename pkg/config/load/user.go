// Package load knows how to load application and repository data
package load

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/user"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

// ErrOnUserDataNotFound returns a handler that errors
// when the app configuration cannot be loaded
func ErrOnUserDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		f, _ := c.GetUserDataPath()

		return &DataNotFoundErr{
			err: fmt.Errorf("failed to load app configuration from: %s", f),
		}
	}
}

// CreateOnUserDataNotFound will start an interactive survey
// that allows the end user to configure okctl for their
// use
func CreateOnUserDataNotFound() DataNotFoundFn {
	return func(c *config.Config) error {
		var err error

		err = PromptContinue(
			"First time using okctl? Start guided configuration?",
			"user aborted interactive configuration",
		)
		if err != nil {
			return err
		}

		userDataPath, err := c.GetUserDataPath()
		if err != nil {
			return err
		}

		err = PromptContinue(
			fmt.Sprintf("Configuration will be written to: %s. Continue?", userDataPath),
			"user aborted interactive configuration",
		)
		if err != nil {
			return err
		}

		data, err := user.New().Survey()
		if err != nil {
			return errors.Wrap(err, "failed to get interactive user data")
		}

		c.UserData = data

		err = c.WriteCurrentUserData()
		if err != nil {
			return errors.Wrap(err, "failed to write current app data")
		}

		c.Logger.WithFields(logrus.Fields{
			"configuration_file": userDataPath,
		}).Info("cli configuration completed")

		return nil
	}
}

// LoaderFn defines an interface for loading configuration
type LoaderFn func(cfg *config.Config, v *viper.Viper) error

// UserDataFromFlagsEnvConfigDefaults returns the default behavior for loading
// application state
func UserDataFromFlagsEnvConfigDefaults(cmd *cobra.Command, notFoundFn DataNotFoundFn) config.DataLoaderFn {
	return buildUserDataLoader(
		loadDefaultUserData,
		loadStoredUserData(notFoundFn),
		loadEnvUserData,
		loadFlagsUserData(cmd),
	)
}

func loadDefaultUserData(_ *config.Config, v *viper.Viper) error {
	b, err := yaml.Marshal(user.New())
	if err != nil {
		return err
	}

	defaultUserData := bytes.NewReader(b)

	v.SetConfigType(config.DefaultConfigType)

	err = v.MergeConfig(defaultUserData)
	if err != nil {
		return err
	}

	return nil
}

func loadStoredUserData(notFoundFn DataNotFoundFn) LoaderFn {
	return func(cfg *config.Config, v *viper.Viper) error {
		configPath, err := cfg.GetUserDataDir()
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

func loadEnvUserData(_ *config.Config, v *viper.Viper) error {
	v.AutomaticEnv()
	v.SetEnvPrefix(config.EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}

func loadFlagsUserData(cmd *cobra.Command) LoaderFn {
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

func updateKnownBinaries(cfg *config.Config) {
	candidates := user.KnownBinaries()

	var update []user.Binary

	for _, candidate := range candidates {
		found := false

		for _, existing := range cfg.UserData.Binaries {
			if candidate.Name == existing.Name && candidate.Version == existing.Version {
				found = true
				break
			}
		}

		if !found {
			update = append(update, candidate)
		}
	}

	cfg.UserData.Binaries = append(cfg.UserData.Binaries, update...)
}

func buildUserDataLoader(loaders ...LoaderFn) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		var err error

		cfg.UserData = &user.Data{}

		v := viper.New()
		v.SetFs(cfg.FileSystem.Fs)

		for _, fn := range loaders {
			err = fn(cfg, v)
			if err != nil {
				return err
			}
		}

		err = v.Unmarshal(cfg.UserData)
		if err != nil {
			return err
		}

		updateKnownBinaries(cfg)

		return nil
	}
}
