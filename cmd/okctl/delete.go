package main

import (
	"errors"
	"fmt"

	"github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/context"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/api/core/cleanup"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	deleteClusterArgs    = 0
	deleteHostedZoneFlag = "i-know-what-i-am-doing-delete-hosted-zone-and-records"
)

func buildDeleteCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete commands",
	}

	deleteClusterCommand := buildDeleteClusterCommand(o)
	cmd.AddCommand(deleteClusterCommand)
	cmd.AddCommand(buildDeletePostgresCommand(o))
	deleteClusterCommand.Flags().String(deleteHostedZoneFlag, "false", "Delete hosted zone")

	return cmd
}

// DeleteClusterOpts contains the required inputs
type DeleteClusterOpts struct {
	AWSCredentialsType    string
	GithubCredentialsType string

	DisableSpinner bool
	Confirm        bool

	Region       string
	AWSAccountID string
	ClusterName  string
}

// Validate the inputs
func (o *DeleteClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: gocyclo, funlen, gocognit
func buildDeleteClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &DeleteClusterOpts{}

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Long: `Delete all resources related to an EKS cluster,
including VPC, this is a highly destructive operation.`,
		Args: cobra.ExactArgs(deleteClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			o.AWSCredentialsType = opts.AWSCredentialsType
			o.GithubCredentialsType = opts.GithubCredentialsType

			err := o.Initialise()
			if err != nil {
				return fmt.Errorf("initialising: %w", err)
			}

			opts.Region = o.Declaration.Metadata.Region
			opts.AWSAccountID = o.Declaration.Metadata.AccountID
			opts.ClusterName = o.Declaration.Metadata.Name

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
				ClusterName:  opts.ClusterName,
			}

			delzones, _ := cmd.Flags().GetString(deleteHostedZoneFlag)

			ready, err := checkifReady(id.ClusterName, o, opts.Confirm)
			if err != nil || !ready {
				return err
			}

			userDir, err := o.GetUserDataDir()
			if err != nil {
				return fmt.Errorf("getting user data: %w", err)
			}

			services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
			if err != nil {
				return fmt.Errorf("creating client services: %w", err)
			}

			formatErr := o.ErrorFormatter(fmt.Sprintf("delete cluster %s", opts.ClusterName), userDir)

			domain, err := services.Domain.GetPrimaryHostedZone(o.Ctx)
			if err != nil {
				return formatErr(err)
			}

			err = services.Monitoring.DeleteTempo(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = services.Monitoring.DeletePromtail(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = services.Monitoring.DeleteLoki(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = services.Monitoring.DeleteKubePromStack(o.Ctx, client.DeleteKubePromStackOpts{
				ID:     id,
				Domain: domain.Domain,
			})
			if err != nil {
				return formatErr(err)
			}

			err = services.ArgoCD.DeleteArgoCD(o.Ctx, client.DeleteArgoCDOpts{
				ID: id,
			})
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

			err = services.AWSLoadBalancerControllerService.DeleteAWSLoadBalancerController(o.Ctx, id)
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

			err = services.Autoscaler.DeleteAutoscaler(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			err = services.Blockstorage.DeleteBlockstorage(o.Ctx, id)
			if err != nil {
				return formatErr(err)
			}

			vpc, err := services.Vpc.GetVPC(o.Ctx, id)
			if err != nil && !errors.Is(err, storm.ErrNotFound) {
				return formatErr(err)
			}

			// Even though we could not retrieve the vpc, we will
			// try to delete the cluster and get rid of as much as
			// possible
			if vpc != nil {
				err = cleanup.DeleteDanglingALBs(o.CloudProvider, vpc.VpcID)
				if err != nil {
					return formatErr(err)
				}

				err = cleanup.DeleteDanglingSecurityGroups(o.CloudProvider, vpc.VpcID)
				if err != nil {
					return formatErr(err)
				}
			}

			err = services.Cluster.DeleteCluster(o.Ctx, client.ClusterDeleteOpts{
				ID:                 id,
				FargateProfileName: "fp-default",
			})
			if err != nil {
				return formatErr(err)
			}

			err = services.Vpc.DeleteVpc(o.Ctx, client.DeleteVpcOpts{
				ID: id,
			})
			if err != nil {
				return formatErr(err)
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.AWSCredentialsType,
		"aws-credentials-type",
		"a",
		context.AWSCredentialsTypeSAML,
		fmt.Sprintf(
			"The form of authentication to use for AWS. Possible values: [%s,%s]",
			context.AWSCredentialsTypeSAML,
			context.AWSCredentialsTypeAccessKey,
		),
	)
	flags.StringVarP(
		&opts.GithubCredentialsType,
		"github-credentials-type",
		"g",
		context.GithubCredentialsTypeDeviceAuthentication,
		fmt.Sprintf(
			"The form of authentication to use for Github. Possible values: [%s,%s]",
			context.GithubCredentialsTypeDeviceAuthentication,
			context.GithubCredentialsTypeToken,
		),
	)

	flags.BoolVar(&opts.DisableSpinner, "no-spinner", false, "disables progress spinner")
	flags.BoolVarP(&opts.Confirm, "confirm", "y", false, "confirm all choices")

	return cmd
}

func checkifReady(clusterName string, o *okctl.Okctl, preconfirmed bool) (bool, error) {
	if preconfirmed {
		return true, nil
	}

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
