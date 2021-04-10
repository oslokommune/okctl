package load

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/config/constant"

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
		baseDir, err := cfg.GetRepoDir()
		if err != nil {
			return err
		}

		exists, err := cfg.FileSystem.Exists(path.Join(baseDir, configFile))
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("couldn't find config file: %s", path.Join(baseDir, configFile))
		}

		_, err = store.NewFileSystem(baseDir, cfg.FileSystem).
			GetStruct(constant.DefaultRepositoryStateFile, cfg.Declaration, store.FromYAML()).
			Do()
		if err != nil {
			return err
		}

		return nil
	}
}
