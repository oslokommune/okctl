package main

import (
	"os"

	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
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

			{
				repoDataNotFound := load.CreateOnRepoDataNotFound()
				if o.NoInput {
					repoDataNotFound = load.ErrOnRepoDataNotFound()
				}
				o.RepoDataLoader = load.RepoDataFromConfigFile(cmd, repoDataNotFound)
			}

			appDataNotFound := load.CreateOnAppDataNotFound()
			if o.NoInput {
				appDataNotFound = load.ErrOnAppDataNotFound()
			}
			o.AppDataLoader = load.AppDataFromFlagsThenEnvVarsThenConfigFile(cmd, appDataNotFound)

			err = o.LoadAppData()
			if err != nil {
				return err
			}

			err = o.LoadRepoData()
			if err != nil {
				return err
			}

			o.Out = cmd.OutOrStdout()
			o.Err = cmd.OutOrStderr()

			return nil
		},
	}

	cmd.AddCommand(buildCreateCommand(o))

	return cmd
}
