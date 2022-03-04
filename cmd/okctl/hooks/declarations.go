package hooks

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/go-git/go-git/v5"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// LoadClusterDeclaration loads a cluster declaration path from disk or stdin
// and stores it in o.Declaration
func LoadClusterDeclaration(o *okctl.Okctl, declarationPath *string) RunEer {
	return func(cmd *cobra.Command, _ []string) error {
		clusterDeclaration, err := LoadClusterDeclarationFromPath(o, declarationPath)
		if err != nil {
			if errors.Is(err, git.ErrRepositoryNotExists) {
				return fmt.Errorf("okctl needs to be run inside a Git repository (okctl outputs " +
					"various configuration files that will bwe stored here)")
			}

			return fmt.Errorf("loading repository data: %w", err)
		}

		o.Declaration = clusterDeclaration

		return nil
	}
}

// LoadClusterDeclarationFromPath loads a cluster declaration path from disk or stdin
func LoadClusterDeclarationFromPath(o *okctl.Okctl, path *string) (*v1alpha1.Cluster, error) {
	emptyCluster := &v1alpha1.Cluster{}

	if len(*path) == 0 {
		return emptyCluster, fmt.Errorf("declaration must be provided")
	}

	var err error

	var pathValue string
	pathValue = *path

	// "-":  user wants to read from stdin
	if pathValue != "-" {
		pathValue, err = filepath.Abs(*path)
		if err != nil {
			return emptyCluster, fmt.Errorf("converting declaration path to absolute path: %w", err)
		}
	}

	clusterDeclaration, err := commands.InferClusterFromStdinOrFile(o.In, pathValue)
	if err != nil {
		return emptyCluster, fmt.Errorf("inferring cluster: %w", err)
	}

	err = clusterDeclaration.Validate()
	if err != nil {
		return emptyCluster, fmt.Errorf("validating cluster declaration: %w", err)
	}

	return clusterDeclaration, nil
}
