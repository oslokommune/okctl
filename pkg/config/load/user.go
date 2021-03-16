// Package load knows how to load application and repository data
package load

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/ask"
	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/oslokommune/okctl/pkg/config"
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

const localMsg = `
Since you haven't used okctl before, we will now setup a directory
where we can store state specific to your user, e.g., binaries, logs,
credentials, and some other configuration:

(config): %s
(data)  : %s

We will also ask some further questions about the setup of your
user. One of these questions is the use of a keyring, encrypted
file or similar.

We use this functionality to store your password and a Github
token securely, so you don't have to enter your credentials
all the time, if you so choose.

`

// CreateOnUserDataNotFoundWithNoInput initializes a user without user input
func CreateOnUserDataNotFoundWithNoInput() DataNotFoundFn {
	return func(c *config.Config) (err error) {
		userDir, err := c.GetUserDataDir()
		if err != nil {
			return fmt.Errorf("getting user data dir: %w", err)
		}

		userDataPath, err := c.GetUserDataPath()
		if err != nil {
			return fmt.Errorf("getting user data path: %w", err)
		}

		data := state.NewUser()

		c.UserState = data

		_, err = store.NewFileSystem(userDir, c.FileSystem).
			StoreStruct(constant.DefaultConfig, c.UserState, store.ToYAML()).
			Do()
		if err != nil {
			return err
		}

		c.Logger.WithFields(logrus.Fields{
			"configuration_file": userDataPath,
		}).Info("cli configuration completed")

		return nil
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

		userDir, err := c.GetUserDataDir()
		if err != nil {
			return fmt.Errorf("getting user data dir: %w", err)
		}

		userDataPath, err := c.GetUserDataPath()
		if err != nil {
			return fmt.Errorf("getting user data path: %w", err)
		}

		_, err = fmt.Fprintf(c.Err, localMsg, userDataPath, userDir)
		if err != nil {
			return err
		}

		err = PromptContinue("Do you want to continue?", "user aborted interactive configuration")
		if err != nil {
			return err
		}

		data := state.NewUser()

		username, err := ask.New().Username()
		if err != nil {
			return errors.Wrap(err, "failed to get interactive user data")
		}

		data.User.Username = username

		c.UserState = data

		_, err = store.NewFileSystem(userDir, c.FileSystem).
			StoreStruct(constant.DefaultConfig, c.UserState, store.ToYAML()).
			Do()
		if err != nil {
			return err
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
	b, err := yaml.Marshal(state.NewUser())
	if err != nil {
		return err
	}

	defaultUserData := bytes.NewReader(b)

	v.SetConfigType(constant.DefaultConfigType)

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
		v.SetConfigName(constant.DefaultConfigName)

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
	v.SetEnvPrefix(constant.EnvPrefix)
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
	candidates := state.KnownBinaries()

	var update []state.Binary

	for _, candidate := range candidates {
		found := false

		for _, existing := range cfg.UserState.Binaries {
			if candidate.Name == existing.Name && candidate.Version == existing.Version {
				found = true
				break
			}
		}

		if !found {
			update = append(update, candidate)
		}
	}

	cfg.UserState.Binaries = append(cfg.UserState.Binaries, update...)
}

func buildUserDataLoader(loaders ...LoaderFn) config.DataLoaderFn {
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

		loadedState := &state.User{}

		// This unmarshal wipes the user
		// name.
		err = v.Unmarshal(loadedState)
		if err != nil {
			return err
		}

		// So we do this weird thing on first run, for now.
		if len(cfg.UserState.User.Username) > 0 {
			loadedState.User.Username = cfg.UserState.User.Username
		}

		cfg.UserState = loadedState

		updateKnownBinaries(cfg)

		return nil
	}
}
