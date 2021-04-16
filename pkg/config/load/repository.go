package load

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/spf13/viper"
)

// RepoDataFromConfigFile defines the default behavior for loading configuration data
func RepoDataFromConfigFile(configFile string) config.DataLoaderFn {
	return buildRepoDataLoader(configFile, nil)
}

func buildRepoDataLoader(configFile string, _ func(v *viper.Viper)) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		repoDir, err := cfg.GetRepoDir()
		if err != nil {
			return err
		}

		cfgFile := path.Join(repoDir, configFile)

		exists, err := cfg.FileSystem.Exists(cfgFile)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("couldn't find config file: %s", cfgFile)
		}

		declaration := v1alpha1.NewCluster()

		_, err = store.NewFileSystem(repoDir, cfg.FileSystem).
			GetStruct(configFile, &declaration, store.FromYAML()).
			Do()
		if err != nil {
			return err
		}

		err = declaration.Validate()
		if err != nil {
			return err
		}

		cfg.Declaration = &declaration

		return nil
	}
}
