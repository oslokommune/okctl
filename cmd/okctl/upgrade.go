package main

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/upgrade"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func buildUpgradeCommand(o *okctl.Okctl) *cobra.Command {
	var upgrader upgrade.Upgrader

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades okctl managed resources to the current version of okctl",
		Long: `Runs a series of upgrade migrations to upgrade resources made by okctl
to the current version of okctl. Example of such resources are helm charts, okctl cluster and application declarations,
binaries used by okctl (kubectl, etc), and internal state.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := o.Initialise()
			if err != nil {
				return err
			}

			handlers := o.StateHandlers(o.StateNodes())
			services, err := o.ClientServices(handlers)
			if err != nil {
				return err
			}

			out := ioutil.Discard
			if o.Debug {
				out = o.Err
			}

			userDataDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			fetcherOpts := upgrade.FetcherOpts{
				Host:  o.Host(),
				Store: storage.NewFileSystemStorage(userDataDir),
			}

			upgrader = upgrade.NewUpgrader(
				o.Logger,
				out,
				services.Github,
				upgrade.NewGithubReleaseParser(upgrade.NewChecksumDownloader()),
				fetcherOpts,
			)

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			err := upgrader.Run()
			if err != nil {
				return fmt.Errorf("upgrading: %w", err)
			}

			return nil
		},
	}

	return cmd
}
