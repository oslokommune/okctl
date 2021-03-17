package main

import (
	"fmt"

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
		Use:   "postgres ENV APPLICATION_NAME",
		Short: "Attach to the given postgres database",
		Args:  cobra.ExactArgs(2), // nolint: gomnd
		PreRunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]
			applicationName := args[1]

			err := o.InitialiseWithOnlyEnv(environment)
			if err != nil {
				return err
			}

			meta := o.RepoStateWithEnv.GetMetadata()
			cluster := o.RepoStateWithEnv.GetCluster()
			db := o.RepoStateWithEnv.GetDatabase(applicationName)

			opts.ID.Environment = environment
			opts.ID.AWSAccountID = cluster.AWSAccountID
			opts.ID.Repository = meta.Name
			opts.ID.Region = meta.Region
			opts.ID.ClusterName = cluster.Name
			opts.ApplicationName = applicationName
			opts.Namespace = db.Namespace
			opts.ConfigMapName = db.DatabaseConfigMapName
			opts.SecretName = db.AdminSecretName
			opts.SecurityGroup = db.SecurityGroupID

			err = opts.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			kubeConfigStore, err := o.KubeConfigStore()
			if err != nil {
				return err
			}

			cfg, err := kubeConfigStore.GetKubeConfig()
			if err != nil {
				return err
			}

			clientSet, config, err := kube.NewFromKubeConfig(cfg.Path).Get()
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

			err = client.Watch(pod)
			if err != nil {
				return err
			}

			err = client.Attach()
			if err != nil {
				return err
			}

			err = client.Delete()
			if err != nil {
				return err
			}

			return policyClient.Delete()
		},
	}

	return cmd
}
