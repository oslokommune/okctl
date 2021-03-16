package main

import (
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const scaffoldClusterArgumentQuantity = 0

func buildScaffoldClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := commands.ScaffoldClusterOpts{}

	cmd := &cobra.Command{
		Use:     "cluster",
		Example: exampleUsage,
		Short:   "Scaffold cluster resource template",
		Long:    "Scaffolds a cluster resource which can be used to control cluster resources",
		Args:    cobra.ExactArgs(scaffoldClusterArgumentQuantity),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ScaffoldClusterDeclaration(o.Out, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "my-product-name", usageName)
	flags.StringVarP(&opts.Environment, "environment", "e", "development", usageEnvironment)
	flags.StringVarP(&opts.Organization, "github-organization", "o", "oslokommune", usageOrganization)
	flags.StringVarP(&opts.RepositoryName, "repository-name", "r", "my_iac_repo_name", usageRepository)
	flags.StringVarP(&opts.AWSAccountID, "aws-account-id", "i", "123456789123", usageAWSAccountID)
	flags.StringVarP(&opts.OutputDirectory, "output-directory", "d", constant.DefaultOutputDirectory, usageOutputDirectory)

	return cmd
}

const (
	usageName            = `name of the cluster`
	usageEnvironment     = `environment for the cluster, for example dev or production`
	usageAWSAccountID    = `aws account where the resources provisioned by okctl should reside`
	usageOrganization    = `organization that owns the infrastructure-as-code repository`
	usageRepository      = `name of the repository that will contain infrastructure-as-code`
	usageOutputDirectory = `name of the directory where okctl will place all infrastructure files`
	exampleUsage         = `okctl scaffold cluster utviklerportalen production > cluster.yaml`
)
