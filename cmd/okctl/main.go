package main

import (
	"os"

	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func binariesProvider(o *okctl.Okctl) error {
	appDataDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(appDataDir)

	stagers, err := binaries.New(o.AppData.Host, store).FromConfig(true, o.AppData.Binaries)
	if err != nil {
		return err
	}

	o.BinariesProvider = stagers

	return nil
}

func repoDataLoader(o *okctl.Okctl, cmd *cobra.Command) error {
	repoDataNotFound := load.CreateOnRepoDataNotFound()

	if o.NoInput {
		repoDataNotFound = load.ErrOnRepoDataNotFound()
	}

	o.RepoDataLoader = load.RepoDataFromConfigFile(cmd, repoDataNotFound)

	return o.LoadRepoData()
}

func appDataLoader(o *okctl.Okctl, cmd *cobra.Command) error {
	appDataNotFound := load.CreateOnAppDataNotFound()

	if o.NoInput {
		appDataNotFound = load.ErrOnAppDataNotFound()
	}

	o.AppDataLoader = load.AppDataFromFlagsEnvConfigDefaults(cmd, appDataNotFound)

	return o.LoadAppData()
}

func buildRootCommand() *cobra.Command {
	o := okctl.New()

	var cmd = &cobra.Command{
		Use:   "okctl",
		Short: "Opinionated and effortless infrastructure and application management",
		Long: `A highly opinionated CLI for creating a Kubernetes cluster in AWS with
a set of applications that ensure tighter integration between AWS and
Kubernetes, e.g., aws-alb-ingress-controller, external-secrets, etc.

Also comes pre-configured with ArgoCD for managing deployments, etc.
We also use the prometheus-operator for ensuring metrics and logs are
being captured. Together with slack and slick.`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			err = appDataLoader(o, cmd)
			if err != nil {
				return err
			}

			err = repoDataLoader(o, cmd)
			if err != nil {
				return err
			}

			err = binariesProvider(o)
			if err != nil {
				return err
			}

			o.Out = cmd.OutOrStdout()
			o.Err = cmd.OutOrStderr()

			return nil
		},
	}

	cmd.AddCommand(buildCreateCommand(o))
	cmd.AddCommand(buildDeleteCommand(o))

	return cmd
}
