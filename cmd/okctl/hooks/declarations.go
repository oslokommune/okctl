package hooks

import (
	"fmt"

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
