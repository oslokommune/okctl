package main

import (
	"fmt"
	"regexp"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/api/core/cleanup"

	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl/pkg/spinner"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	deleteClusterArgs    = 1
	deleteHostedZoneFlag = "i-know-what-i-am-doing-delete-hosted-zone-and-records"
)

func buildDeleteCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete commands",
	}

	deleteClusterCommand := buildDeleteClusterCommand(o)
	cmd.AddCommand(deleteClusterCommand)
	cmd.AddCommand(buildDeleteTestClusterCommand(o))
	deleteClusterCommand.Flags().String(deleteHostedZoneFlag, "false", "Delete hosted zone")

	return cmd
}

// DeleteClusterOpts contains the required inputs
type DeleteClusterOpts struct {
	Region       string
	AWSAccountID string
	Environment  string
	Repository   string
	ClusterName  string
}

// Validate the inputs
func (o *DeleteClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: gocyclo, funlen, gocognit
func buildDeleteClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &DeleteClusterOpts{}

	cmd := &cobra.Command{
		Use:   "cluster [env]",
		Short: "Delete a cluster",
		Long: `Delete all resources related to an EKS cluster,
including VPC, this is a highly destructive operation.`,
		Args: cobra.ExactArgs(deleteClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]

			err := validation.Validate(
				&environment,
				validation.Required,
				validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("the environment must consist of 3-64 characters (a-z, A-Z)"),
			)
			if err != nil {
				return err
			}

			err = o.InitialiseWithOnlyEnv(environment)
			if err != nil {
				return err
			}

			meta := o.RepoStateWithEnv.GetMetadata()
			cluster := o.RepoStateWithEnv.GetCluster()

			opts.Repository = meta.Name
			opts.Region = meta.Region
			opts.AWSAccountID = cluster.AWSAccountID
			opts.Environment = cluster.Environment
			opts.ClusterName = cluster.Name

			err = opts.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			id := api.ID{
				Region:       opts.Region,
				AWSAccountID: opts.AWSAccountID,
				Environment:  opts.Environment,
				Repository:   opts.Repository,
				ClusterName:  opts.ClusterName,
			}

			delzones, _ := cmd.Flags().GetString(deleteHostedZoneFlag)

			ready, err := checkifReady(id.ClusterName, o)
			if err != nil || !ready {
				return err
			}

			userDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			spin, err := spinner.New("deleting", o.Err)
			if err != nil {
				return err
			}

			services, err := o.ClientServices(spin)
			if err != nil {
				return err
			}

			formatErr := o.ErrorFormatter(fmt.Sprintf("delete cluster %s", opts.Environment), userDir)

			cluster := o.RepoStateWithEnv.GetCluster()

			vpc, err := services.Vpc.GetVPC(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			if delzones == "true" {
				err = services.Domain.DeletePrimaryHostedZone(o.Ctx, client.DeletePrimaryHostedZoneOpts{
					ID: id,
				})
				if err != nil {
					return formatErr(err)
				}
			}

			err = services.IdentityManager.DeleteIdentityPool(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = services.ALBIngressController.DeleteALBIngressController(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = services.ExternalSecrets.DeleteExternalSecrets(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = services.ExternalDNS.DeleteExternalDNS(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = cleanup.DeleteDanglingALBs(o.CloudProvider, vpc.VpcID)
			if err != nil {
				return formatErr(err)
			}

			err = services.Cluster.DeleteCluster(o.Ctx, api.ClusterDeleteOpts{
				ID:                 id,
				FargateProfileName: "fp-default",
			})
			if err != nil {
				return formatErr(err)
			}

			err = services.Parameter.DeleteAllsecrets(o.Ctx, cluster)
			if err != nil {
				return formatErr(err)
			}

			err = cleanup.DeleteDanglingSecurityGroups(o.CloudProvider, vpc.VpcID)
			if err != nil {
				return formatErr(err)
			}

			err = services.Vpc.DeleteVpc(o.Ctx, api.DeleteVpcOpts{
				ID: id,
			})
			if err != nil {
				return formatErr(err)
			}

			return nil
		},
	}

	return cmd
}

func checkifReady(clusterName string, o *okctl.Okctl) (bool, error) {
	ready := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("This will delete %s and all assosicated resources, are you sure you want to continue?", clusterName),
	}

	err := survey.AskOne(prompt, &ready)
	if err != nil {
		return false, err
	}

	if !ready {
		_, err = fmt.Fprintf(o.Err, "user wasn't ready to continue, aborting.")
		if err != nil {
			return false, err
		}

		return false, err
	}

	return true, nil
}
