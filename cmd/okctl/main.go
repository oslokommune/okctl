package main

import (
	"os"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
	repoCfg := &config.RepoConfig{}

	appCfg := &config.AppConfig{}

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
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			var err error

			appCfg, err = config.LoadAppCfg()
			if err != nil {
				return err
			}

			switch err.(type) {
			case *config.AppCfgNotFoundErr:
				// Here we need to perform a survey, asking the
				// user all kinds of questions.
				appCfg, err = config.NewAppCfg()
				if err != nil {
					return err
				}
			default:
				return err
			}

			repoCfg, err = config.LoadRepoCfg()
			if err != nil {
				return err
			}

			return err
		},
	}

	cmd.AddCommand(buildCreateCommand(appCfg, repoCfg))
	cmd.AddCommand(buildLoginCommand(appCfg, repoCfg))

	return cmd
}
