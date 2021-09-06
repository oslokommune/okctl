package main

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/upgrade/clusterversioner"
	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversioner"
	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/upgrade"
	"github.com/spf13/cobra"
)

type upgradeFlags struct {
	confirm bool
}

//nolint:funlen
func buildUpgradeCommand(o *okctl.Okctl) *cobra.Command {
	flags := upgradeFlags{}

	var upgrader upgrade.Upgrader

	var originalClusterVersioner originalclusterversioner.Versioner

	var clusterVersioner clusterversioner.Versioner

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

			clusterVersioner = clusterversioner.New(
				out,
				api.ID{
					Region:       o.Declaration.Metadata.Region,
					AWSAccountID: o.Declaration.Metadata.AccountID,
					ClusterName:  o.Declaration.Metadata.Name,
				},
				stateHandlers.Upgrade,
			)

			originalClusterVersioner = originalclusterversioner.New(
				api.ID{
					Region:       o.Declaration.Metadata.Region,
					AWSAccountID: o.Declaration.Metadata.AccountID,
					ClusterName:  o.Declaration.Metadata.Name,
				},
				stateHandlers.Upgrade,
				stateHandlers.Cluster,
			)

			err = clusterVersioner.ValidateBinaryVsClusterVersion(version.GetVersionInfo().Version)
			if err != nil {
				return fmt.Errorf(commands.ValidateBinaryVsClusterVersionError, err)
			}

			err = originalClusterVersioner.SaveOriginalClusterVersionIfNotExists()
			if err != nil {
				return fmt.Errorf(originalclusterversioner.SaveErrorMessage, err)
			}

			originalClusterVersion, err := stateHandlers.Upgrade.GetOriginalClusterVersion()
			if err != nil {
				return fmt.Errorf("getting original okctl version: %w", err)
			}

			upgrader = upgrade.New(upgrade.Opts{
				Debug:                    o.Debug,
				Logger:                   o.Logger,
				Out:                      out,
				AutoConfirmPrompt:        flags.confirm,
				RepositoryDirectory:      repoDir,
				GithubService:            services.Github,
				ChecksumDownloader:       upgrade.NewChecksumDownloader(),
				ClusterVersioner:         clusterVersioner,
				OriginalClusterVersioner: originalClusterVersioner,
				FetcherOpts:              fetcherOpts,
				OkctlVersion:             version.GetVersionInfo().Version,
				OriginalClusterVersion:   originalClusterVersion.Value,
				State:                    stateHandlers.Upgrade,
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
		Hidden: true,
	}

	cmd.PersistentFlags().BoolVarP(&flags.confirm, "confirm", "y", false, "Skip the confirmation prompt")

	return cmd
}
