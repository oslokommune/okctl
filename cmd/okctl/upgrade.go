package main

import (
	"fmt"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/upgrade/clusterversion"
	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversion"

	"github.com/oslokommune/okctl/pkg/upgrade/survey"

	"github.com/oslokommune/okctl/pkg/api"
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

	var originalClusterVersioner originalclusterversion.Versioner

	var clusterVersioner clusterversion.Versioner

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades okctl managed resources to the current version of okctl",
		Long: `Runs a series of upgrade migrations to upgrade resources made by okctl
to the current version of okctl. Example of such resources are helm charts, okctl cluster and application declarations,
binaries used by okctl (kubectl, etc), and internal state.`,
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionUpgrade),
			hooks.InitializeOkctl(o),
			hooks.AcquireStateLock(o),
			hooks.DownloadState(o, true),
			hooks.VerifyClusterExistsInState(o),
			func(cmd *cobra.Command, args []string) error {
				okctlEnvironment, err := commands.GetOkctlEnvironment(o, declarationPath)
				if err != nil {
					return err
				}

				upgradeBinaryEnvVars, err := commands.GetVenvEnvVars(okctlEnvironment)
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

				repoDir, err := o.GetRepoDir()
				if err != nil {
					return err
				}

				fetcherOpts := upgrade.FetcherOpts{
					Host:  o.Host(),
					Store: storage.NewFileSystemStorage(userDataDir),
				}

				clusterID := api.ID{
					Region:       o.Declaration.Metadata.Region,
					AWSAccountID: o.Declaration.Metadata.AccountID,
					ClusterName:  o.Declaration.Metadata.Name,
				}

				clusterVersioner = clusterversion.New(
					out,
					clusterID,
					stateHandlers.Upgrade,
				)

				surveyor := survey.NewTerminalSurveyor(out, flags.confirm)

				originalClusterVersioner = originalclusterversion.New(clusterID, stateHandlers.Upgrade, stateHandlers.Cluster)

				upgrader, err = upgrade.New(upgrade.Opts{
					Debug:                      o.Debug,
					AutoConfirm:                flags.confirm,
					Logger:                     o.Logger,
					Out:                        out,
					RepositoryDirectory:        repoDir,
					GithubService:              services.Github,
					ChecksumDownloader:         upgrade.NewChecksumDownloader(),
					ClusterVersioner:           clusterVersioner,
					OriginalClusterVersioner:   originalClusterVersioner,
					Surveyor:                   surveyor,
					FetcherOpts:                fetcherOpts,
					OkctlVersion:               version.GetVersionInfo().Version,
					State:                      stateHandlers.Upgrade,
					ClusterID:                  clusterID,
					BinaryEnvironmentVariables: upgradeBinaryEnvVars,
				})
				if err != nil {
					return fmt.Errorf("creating upgrader: %w", err)
				}

				return nil
			},
		),
		RunE: func(_ *cobra.Command, args []string) error {
			err := upgrader.Run()
			if err != nil {
				return fmt.Errorf("upgrading: %w", err)
			}
			return nil
		},
		PostRunE: hooks.RunECombinator(
			hooks.UploadState(o),
			hooks.ClearLocalState(o),
			hooks.ReleaseStateLock(o),
			hooks.EmitEndCommandExecutionEvent(metrics.ActionUpgrade),
		),
	}

	cmd.PersistentFlags().BoolVarP(
		&flags.confirm, "confirm", "y", false, "Skip confirmation prompts")

	return cmd
}
