package main

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/oslokommune/okctl/cmd/login"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	short = "Opinionated and effortless infrastructure and application management"
	long  = `A highly opinionated CLI for creating a Kubernetes cluster in AWS with
a set of applications that ensure tighter integration between AWS and
Kubernetes, e.g., aws-alb-ingress-controller, external-secrets, etc.

Also comes pre-configured with ArgoCD for managing deployments, etc.
We also use the prometheus-operator for ensuring metrics and logs are
being captured. Together with slack and slick.`
)

func main() {
	cmd := BuildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func BuildRootCommand() *cobra.Command {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("failed to get home directory: %s", err)
		os.Exit(1)
	}

	cfg, err := config.LoadUserConfiguration(home)
	if err != nil {
		fmt.Printf("failed to load user configuration: %s", err)
		os.Exit(1)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(logger.Writer())

	var cmd = &cobra.Command{
		Use:   "okctl",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(login.BuildLoginCommand(cfg, logger))

	return cmd
}
