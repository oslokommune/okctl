package main

import (
	"github.com/oslokommune/okctl/cmd/okctl/preruns"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const scaffoldClusterArgumentQuantity = 0

func buildScaffoldClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := commands.ScaffoldClusterOpts{}

	cmd := &cobra.Command{
		Use:     "cluster",
		Example: exampleUsage,
		Short:   ScaffoldClusterShortDescription,
		Long:    ScaffoldClusterLongDescription,
		Args:    cobra.ExactArgs(scaffoldClusterArgumentQuantity),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
		PreRunE: preruns.PreRunECombinator(
			preruns.LoadUserData(o),
			preruns.InitializeMetrics(o),
			func(cmd *cobra.Command, args []string) error {
				metrics.Publish(metrics.Event{
					Category: metrics.CategoryCommandExecution,
					Action:   metrics.ActionScaffoldCluster,
					Label:    metrics.LabelStart,
				})

				return nil
			},
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.ScaffoldClusterDeclaration(o.Out, opts)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			metrics.Publish(metrics.Event{
				Category: metrics.CategoryCommandExecution,
				Action:   metrics.ActionScaffoldCluster,
				Label:    metrics.LabelEnd,
			})

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "my-cluster-name", usageName)
	flags.StringVarP(&opts.Organization, "github-organization", "o", "oslokommune", usageOrganization)
	flags.StringVarP(&opts.RepositoryName, "repository-name", "r", "my_iac_repo_name", usageRepository)
	flags.StringVarP(&opts.AWSAccountID, "aws-account-id", "i", "123456789123", usageAWSAccountID)
	flags.StringVarP(&opts.OutputDirectory, "output-directory", "d", constant.DefaultOutputDirectory, usageOutputDirectory)

	return cmd
}

const (
	usageName            = `name of the cluster`
	usageAWSAccountID    = `aws account where the resources provisioned by okctl should reside`
	usageRepository      = `name of the repository that will contain infrastructure-as-code`
	usageOrganization    = `organization that owns the infrastructure-as-code repository`
	usageOutputDirectory = `name of the directory where okctl will place all infrastructure files`
	exampleUsage         = `okctl scaffold cluster > cluster.yaml`
)
