package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/oslokommune/okctl/pkg/kube/manifests/pgbouncer"

	"github.com/oslokommune/okctl/pkg/kube/manifests/awsnode"

	"github.com/oslokommune/okctl/pkg/kube/manifests/securitygrouppolicy"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type forwardPostgresOpts struct {
	ID              api.ID
	ApplicationName string
	Environment     string
	Namespace       string
	ConfigMapName   string
	SecretName      string
	SecurityGroup   string
	Username        string
	PasswordFile    string
	ListenPort      int32
}

// Validate the inputs
func (o *forwardPostgresOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ApplicationName, validation.Required),
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
		validation.Field(&o.Username, validation.Required),
		validation.Field(&o.PasswordFile, validation.Required),
		validation.Field(&o.ListenPort, validation.Required),
	)
}

// nolint: funlen
func buildForwardPostgres(o *okctl.Okctl) *cobra.Command {
	opts := &forwardPostgresOpts{}

	cmd := &cobra.Command{
		Use:   "postgres",
		Short: "Forward to the given postgres database",
		Args:  cobra.ExactArgs(0), // nolint: gomnd
		PreRunE: func(_ *cobra.Command, _ []string) error {
			err := o.InitialiseWithOnlyEnv(opts.Environment)
			if err != nil {
				return err
			}

			meta := o.RepoStateWithEnv.GetMetadata()
			cluster := o.RepoStateWithEnv.GetCluster()
			db := o.RepoStateWithEnv.GetDatabase(opts.ApplicationName)

			opts.ID.Environment = opts.Environment
			opts.ID.AWSAccountID = cluster.AWSAccountID
			opts.ID.Repository = meta.Name
			opts.ID.Region = meta.Region
			opts.ID.ClusterName = o.RepoStateWithEnv.GetClusterName()
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
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			clientSet, config, err := kube.NewFromEKSCluster(
				opts.ID.ClusterName,
				opts.ID.Region,
				o.CloudProvider,
				o.CredentialsProvider.Aws(),
			).Get()
			if err != nil {
				return err
			}

			// Ensure that ENABLE_POD_ENI is true
			err = awsnode.New(clientSet).EnablePodENI()
			if err != nil {
				return err
			}

			password, err := o.FileSystem.ReadFile(opts.PasswordFile)
			if err != nil {
				return err
			}

			app := fmt.Sprintf("%s-pgbouncer", opts.ApplicationName)

			labels := map[string]string{
				"pgbouncer": app,
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
				perr := policyClient.Delete()
				if perr != nil {
					err = perr
				}
			}()

			client := pgbouncer.New(&pgbouncer.Config{
				Name:                  app,
				Namespace:             opts.Namespace,
				Username:              opts.Username,
				Password:              strings.TrimSpace(string(password)),
				DBParamsSecretName:    opts.SecretName,
				DBParamsConfigmapName: opts.ConfigMapName,
				Labels:                labels,
				ListenPort:            opts.ListenPort,
				In:                    os.Stdin,
				Out:                   os.Stdout,
				Err:                   os.Stderr,
				ClientSet:             clientSet,
				Config:                config,
			})

			err = client.Create()
			if err != nil {
				return err
			}

			defer func() {
				cerr := client.Delete()
				if cerr != nil {
					err = cerr
				}
			}()

			return err
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.Environment,
		"environment",
		"e",
		"",
		"The environment the postgres database was created in",
	)

	flags.StringVarP(&opts.ApplicationName,
		"name",
		"n",
		"",
		"The name of the database to forward to",
	)

	flags.StringVarP(&opts.Username,
		"username",
		"u",
		"",
		"The username you want to set for the session",
	)

	flags.StringVarP(&opts.PasswordFile,
		"password-file",
		"p",
		"",
		"A file containing the password you want to set for the session",
	)

	flags.Int32VarP(&opts.ListenPort,
		"listen-port",
		"l",
		5432,
		"The port on your local machine traffic will be forwarded to",
	)

	return cmd
}
