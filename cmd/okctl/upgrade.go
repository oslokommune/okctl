package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/upgrade/clusterversion"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/upgrade/originalversion"
	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/upgrade"
	"github.com/spf13/cobra"
)

//nolint:funlen
func buildUpgradeCommand(o *okctl.Okctl) *cobra.Command {
	var upgrader upgrade.Upgrader

	var originalVersionSaver originalversion.Saver

	var clusterVersioner clusterversion.ClusterVersioner

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

			out := o.Out
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

			// Cluster version
			clusterVersioner = clusterversion.New(
				out,
				api.ID{
					Region:       o.Declaration.Metadata.Region,
					AWSAccountID: o.Declaration.Metadata.AccountID,
					ClusterName:  o.Declaration.Metadata.Name,
				},
				stateHandlers.Upgrade,
			)

			err = clusterVersioner.ValidateBinaryVsClusterVersion(version.GetVersionInfo().Version)
			if err != nil {
				return fmt.Errorf(commands.ValidateBinaryVsClusterVersionError, err)
			}

			// Original version
			originalVersionSaver, err = originalversion.New(
				api.ID{
					Region:       o.Declaration.Metadata.Region,
					AWSAccountID: o.Declaration.Metadata.AccountID,
					ClusterName:  o.Declaration.Metadata.Name,
				},
				stateHandlers.Upgrade,
				stateHandlers.Cluster,
			)
			if err != nil {
				return fmt.Errorf("creating original version saver: %w", err)
			}

			err = originalVersionSaver.SaveOriginalClusterVersionIfNotExists()
			if err != nil {
				return fmt.Errorf(originalversion.SaveErrorMessage, err)
			}

			originalClusterVersion, err := stateHandlers.Upgrade.GetOriginalClusterVersion()
			if err != nil {
				return fmt.Errorf("getting original okctl version: %w", err)
			}

			upgrader = upgrade.New(upgrade.Opts{
				Debug:                  o.Debug,
				Logger:                 o.Logger,
				Out:                    out,
				RepositoryDirectory:    repoDir,
				GithubService:          services.Github,
				ChecksumDownloader:     upgrade.NewChecksumDownloader(),
				FetcherOpts:            fetcherOpts,
				OkctlVersion:           version.GetVersionInfo().Version,
				OriginalClusterVersion: originalClusterVersion.Value,
				State:                  stateHandlers.Upgrade,
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

			err = clusterVersioner.SaveClusterVersion(version.GetVersionInfo())
			if err != nil {
				return fmt.Errorf(commands.SaveClusterVersionError, err)
			}

			return nil
		},
		Hidden: true,
	}

	return cmd
}
