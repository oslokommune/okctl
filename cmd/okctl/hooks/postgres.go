package hooks

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// ValidatePostgresApplicationName ensures that applicationName matches an existing database
func ValidatePostgresApplicationName(o *okctl.Okctl, applicationName *string) RunEer {
	return func(cmd *cobra.Command, args []string) error {
		databaseName := *applicationName
		if len(databaseName) == 0 {
			return fmt.Errorf("missing database instance name")
		}

		var existingDatabaseNames []string

		dbs, _ := o.StateHandlers(o.StateNodes()).Component.GetPostgresDatabases()
		for _, db := range dbs {
			existingDatabaseNames = append(existingDatabaseNames, db.ApplicationName)
		}

		if !databaseNameExists(databaseName, existingDatabaseNames) {
			return fmt.Errorf(
				"database name '%s' is not valid. The following datbases are available: [%s]",
				databaseName,
				strings.Join(existingDatabaseNames, ", "),
			)
		}

		return nil
	}
}

func databaseNameExists(databaseName string, existingDatabaseNames []string) bool {
	for _, b := range existingDatabaseNames {
		if b == databaseName {
			return true
		}
	}

	return false
}
