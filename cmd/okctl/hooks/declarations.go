package hooks

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/oslokommune/okctl/pkg/config/load"

	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// LoadClusterDeclaration acquires a cluster manifest from either stdin
// or a file and stores it in o.Declaration
func LoadClusterDeclaration(o *okctl.Okctl, path *string) RunEer {
	return func(cmd *cobra.Command, args []string) error {
		declaration, err := commands.InferClusterFromStdinOrFile(o.In, *path)
		if err != nil {
			return fmt.Errorf("inferring cluster: %w", err)
		}

		err = declaration.Validate()
		if err != nil {
			return fmt.Errorf("validating cluster declaration: %w", err)
		}

		o.Declaration = declaration

		return nil
	}
}

// LoadClusterDeclarationPath loads a cluster declaration path from disk
// Shares a lot in common with LoadClusterDeclaration and these should be merged together.
// But: LoadClusterDeclaration also takes into consideration loading from stdin, so we keep
// these two separate for the time being
func LoadClusterDeclarationPath(o *okctl.Okctl, declarationPath *string) RunEer {
	return func(cmd *cobra.Command, _ []string) error {
		var err error

		if len(*declarationPath) == 0 {
			return fmt.Errorf("declaration must be provided")
		}

		tmp, err := filepath.Abs(*declarationPath)
		if err != nil {
			return fmt.Errorf("converting declaration path to absolute path: %w", err)
		}

		declarationPath = &tmp

		err = loadRepoData(o, *declarationPath, cmd)
		if err != nil {
			if errors.Is(err, git.ErrRepositoryNotExists) {
				return fmt.Errorf("okctl needs to be run inside a Git repository (okctl outputs " +
					"various configuration files that will be stored here)")
			}

			return fmt.Errorf("loading repository data: %w", err)
		}

		return nil
	}
}

func loadRepoData(o *okctl.Okctl, declarationPath string, _ *cobra.Command) error {
	o.RepoDataLoader = load.RepoDataFromConfigFile(declarationPath)

	return o.LoadRepoData()
}
