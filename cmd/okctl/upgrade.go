package main

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/upgrade"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades okctl managed resources to the current version of okctl",
		Long: `Runs a series of upgrade migrations to upgrade resources made by okctl
to the current version of okctl. Example of such resources are helm charts, okctl cluster and application declarations,
binaries used by okctl (kubectl, etc), and internal state.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			handlers := o.StateHandlers(o.StateNodes())
			services, err := o.ClientServices(handlers)
			if err != nil {
				return err
			}

			upgrader := upgrade.NewUpgrader(services.Github, services.BinaryService, o.BinariesProvider)

			err = upgrader.Run()
			if err != nil {
				return fmt.Errorf("upgrading: %w", err)
			}

			return nil
		},
	}

	return cmd
}
