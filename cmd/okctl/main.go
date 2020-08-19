package main

import (
	"os"

	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
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
	var outputFormat string

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
				return errors.Wrap(err, "failed to load application data")
			}

			err = repoDataLoader(o, cmd)
			if err != nil {
				return errors.Wrap(err, "failed to load repository data")
			}

			o.Out = cmd.OutOrStdout()
			o.Err = cmd.OutOrStderr()

			o.SetFormat(core.EncodeResponseType(outputFormat))

			return nil
		},
	}

	cmd.AddCommand(buildCreateCommand(o))
	cmd.AddCommand(buildDeleteCommand(o))
	cmd.AddCommand(buildShowCommand(o))
	cmd.AddCommand(buildVersionCommand(o))

	f := cmd.Flags()
	f.StringVarP(&outputFormat, "output", "o", "text",
		"The format of the output returned to the user")

	return cmd
}
