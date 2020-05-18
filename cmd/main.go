package main

import (
	"os"

	"github.com/oslokommune/okctl/cmd/configure"
	"github.com/oslokommune/okctl/cmd/login"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
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
	}

	cmd.AddCommand(login.BuildLoginCommand())
	cmd.AddCommand(configure.BuildConfigureCommand())

	return cmd
}
