package main

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"

	"github.com/oslokommune/okctl/pkg/spinner"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	deleteTestClusterArgs = 1
)

// DeleteTestClusterOpts contains the required inputs
type DeleteTestClusterOpts struct {
	Region       string
	AWSAccountID string
	Environment  string
	Repository   string
	ClusterName  string
}

// Validate the inputs
func (o *DeleteTestClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: funlen gocognit
func buildDeleteTestClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &DeleteTestClusterOpts{}

	cmd := &cobra.Command{
		Use:   "testcluster [env]",
		Short: "Delete a test cluster",
		Long: `Delete all resources related to an EKS test cluster,
including VPC, this is a highly destructive operation.`,
		Args: cobra.ExactArgs(deleteTestClusterArgs),
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
				return errors.E(err, "failed to validate delete testcluster options")
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			id := api.ID{
				Region:       opts.Region,
				AWSAccountID: opts.AWSAccountID,
				Environment:  opts.Environment,
				Repository:   opts.Repository,
				ClusterName:  opts.ClusterName,
			}

			_, err := fmt.Fprintf(o.Err, `
We will now delete your cluster, but first you need
to verify that you have:

Deleted all ingress and service resources that
have created ALBs, you can check what resources
that are running with:

$ kubectl get ingress --all-namespaces
$ kubectl get service --all-namespaces
`)
			if err != nil {
				return err
			}

			userDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			ready := false
			prompt := &survey.Confirm{
				Message: "I confirm that I have removed all ingress and service resources that might create ALBs.",
			}

			err = survey.AskOne(prompt, &ready)
			if err != nil {
				return err
			}

			if !ready {
				_, err = fmt.Fprintf(o.Err, "user wasn't ready to continue, aborting.")
				if err != nil {
					return err
				}

				return nil
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

			// This is taken out, because of possible unintended consequences. The code is kept for now
			/*
				err = services.Domain.DeletePrimaryHostedZone(o.Ctx, o.CloudProvider, client.DeletePrimaryHostedZoneOpts{
					ID: id,
				})
				if err != nil {
					return formatErr(err)
				}
			*/

			err = services.ExternalDNS.DeleteExternalDNS(o.Ctx, id)
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

			err = services.Cluster.DeleteCluster(o.Ctx, api.ClusterDeleteOpts{
				ID:                 id,
				FargateProfileName: "fp-default",
			})
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
