package main

import (
	"fmt"
	"io/ioutil"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/upgrade"
	"github.com/spf13/cobra"
)

//nolint:funlen
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

			stateHandlers := o.StateHandlers(o.StateNodes())

			services, err := o.ClientServices(stateHandlers)
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

			repoDir, err := o.GetHomeDir()
			if err != nil {
				return err
			}

			fetcherOpts := upgrade.FetcherOpts{
				Host:  o.Host(),
				Store: storage.NewFileSystemStorage(userDataDir),
			}

			upgrader = upgrade.New(upgrade.Opts{
				Debug:               o.Debug,
				Logger:              o.Logger,
				Out:                 out,
				RepositoryDirectory: repoDir,
				GithubService:       services.Github,
				ChecksumDownloader:  upgrade.NewChecksumDownloader(),
				FetcherOpts:         fetcherOpts,
				OkctlVersion:        version.String(),
				State:               stateHandlers.Upgrade,
				ClusterID: api.ID{
					Region:       o.Declaration.Metadata.Region,
					AWSAccountID: o.Declaration.Metadata.AccountID,
					ClusterName:  o.Declaration.Metadata.Name,
				},
			})

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
