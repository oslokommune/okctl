package load

import (
	"fmt"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/spf13/viper"
)

// RepoDataFromConfigFile defines the default behavior for loading configuration data
func RepoDataFromConfigFile(declarationPath string) config.DataLoaderFn {
	return buildRepoDataLoader(declarationPath, nil)
}

func buildRepoDataLoader(declarationPath string, _ func(v *viper.Viper)) config.DataLoaderFn {
	return func(cfg *config.Config) error {
		exists, err := cfg.FileSystem.Exists(declarationPath)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("couldn't find config file: %s", declarationPath)
		}

		directory := filepath.Dir(declarationPath)
		base := filepath.Base(declarationPath)

		declaration := v1alpha1.NewCluster()

		_, err = store.NewFileSystem(directory, cfg.FileSystem).
			GetStruct(base, &declaration, store.FromYAML()).
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
