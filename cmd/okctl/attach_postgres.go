package main

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/kube/manifests/awsnode"

	"github.com/oslokommune/okctl/pkg/kube/manifests/securitygrouppolicy"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/oslokommune/okctl/pkg/kube/manifests/psqlclient"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type attachPostgresOpts struct {
	ID              api.ID
	ApplicationName string
	Namespace       string
	ConfigMapName   string
	SecretName      string
	SecurityGroup   string
}

// Validate the inputs
func (o *attachPostgresOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ApplicationName, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
	)
}

// nolint: funlen
func buildAttachPostgres(o *okctl.Okctl) *cobra.Command {
	opts := &attachPostgresOpts{}

	cmd := &cobra.Command{
		Use:   "postgres",
		Short: "Attach to the given postgres database",
		Args:  cobra.ExactArgs(0), // nolint: gomnd
		PreRunE: func(_ *cobra.Command, _ []string) error {
			err := o.Initialise()
			if err != nil {
				return err
			}

			db, err := o.StateHandlers(o.StateNodes()).Component.
				GetPostgresDatabase(cfn.NewStackNamer().
					RDSPostgres(opts.ApplicationName, o.Declaration.Metadata.Name))
			if err != nil {
				return err
			}

			opts.ID.AWSAccountID = o.Declaration.Metadata.AccountID
			opts.ID.Region = o.Declaration.Metadata.Region
			opts.ID.ClusterName = o.Declaration.Metadata.Name
			opts.Namespace = db.Namespace
			opts.ConfigMapName = db.DatabaseConfigMapName
			opts.SecretName = db.AdminSecretName
			opts.SecurityGroup = db.OutgoingSecurityGroupID

			err = opts.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			clientSet, config, err := kube.NewFromEKSCluster(
				opts.ID.ClusterName,
				opts.ID.Region,
				o.CloudProvider,
				o.CredentialsProvider.Aws(),
			).
				Get()
			if err != nil {
				return err
			}

			// Ensure that ENABLE_POD_ENI is true
			err = awsnode.New(clientSet).EnablePodENI()
			if err != nil {
				return err
			}

			app := fmt.Sprintf("%s-psqlclient", opts.ApplicationName)

			labels := map[string]string{
				"psqlclient": app,
			}

			policyClient := securitygrouppolicy.New(
				app,
				opts.Namespace,
				securitygrouppolicy.Manifest(
					app,
					opts.Namespace,
					labels,
					[]string{opts.SecurityGroup},
				),
				config,
			)

			_, err = policyClient.Create()
			if err != nil {
				return err
			}

			defer func() {
				err = policyClient.Delete()
			}()

			client := psqlclient.New(
				app,
				opts.Namespace,
				psqlclient.Manifest(
					app,
					opts.Namespace,
					opts.ConfigMapName,
					opts.SecretName,
					labels,
				),
				clientSet,
				config,
			)

			pod, err := client.Create()
			if err != nil {
				return err
			}

			defer func() {
				err = client.Delete()
			}()

			err = client.Watch(pod)
			if err != nil {
				return err
			}

			err = client.Attach()
			if err != nil {
				return err
			}

			return err
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.ApplicationName,
		"name",
		"n",
		"",
		"The name of the database to attach to",
	)

	return cmd
}
