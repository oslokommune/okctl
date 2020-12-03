package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

const scaffoldClusterArgumentQuantity = 2

type scaffoldClusterOpts struct {
	Name string

	AWSAccountID   string
	Environment    string
	Organization   string
	RepositoryName string
	Team           string
}

func buildScaffoldClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := scaffoldClusterOpts{}

	cmd := &cobra.Command{
		Use:     "cluster CLUSTER_NAME ENVIRONMENT",
		Example: exampleUsage,
		Short:   "Scaffold cluster resource template",
		Long:    "Scaffolds a cluster resource which can be used to control cluster resources",
		Args:    cobra.ExactArgs(scaffoldClusterArgumentQuantity),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Environment = args[1]

			clusterResource := v1alpha1.NewDefaultCluster(
				opts.Name,
				opts.Environment,
				opts.Organization,
				opts.RepositoryName,
				opts.Team,
				opts.AWSAccountID,
			)

			result, err := yaml.Marshal(clusterResource)
			if err != nil {
				return fmt.Errorf("error marshalling yaml: %w", err)
			}

			_, err = o.Out.Write(result)
			if err != nil {
				return fmt.Errorf("error writing cluster: %w", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Organization, "github-organization", "o", "oslokommune", usageOrganization)
	flags.StringVarP(&opts.RepositoryName, "repository-name", "r", "my_iac_repo_name", usageRepository)
	flags.StringVarP(&opts.Team, "github-team", "t", "my_team", usageTeam)
	flags.StringVarP(&opts.AWSAccountID, "aws-account-id", "i", "123456789123", usageAWSAccountID)

	return cmd
}

const (
	usageAWSAccountID = `the aws account where the resources provisioned by okctl should reside`
	usageOrganization = `the organization that owns the infrastructure-as-code repository`
	usageRepository   = `the name of the repository that will contain infrastructure-as-code`
	usageTeam         = `the team that is responsible and has access rights to the infrastructure-as-code repository`
	exampleUsage      = `okctl scaffold cluster utviklerportalen production > cluster.yaml`
)
