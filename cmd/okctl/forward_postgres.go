package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/oslokommune/okctl/cmd/okctl/preruns"
	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/cfn"

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
	DatabaseName    string
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
		Short: ForwardPostgresShortDescription,
		Long:  ForwardPostgresLongDescription,
		Args:  cobra.ExactArgs(0), // nolint: gomnd
		PreRunE: preruns.PreRunECombinator(
			preruns.LoadUserData(o),
			preruns.InitializeMetrics(o),
			func(_ *cobra.Command, _ []string) error {
				metrics.Publish(metrics.Event{
					Category: metrics.CategoryCommandExecution,
					Action:   metrics.ActionForwardPostgres,
					Label:    metrics.LabelStart,
				})

				if len(opts.ApplicationName) == 0 {
					return fmt.Errorf("missing database instance name")
				}

				err := o.Initialise()
				if err != nil {
					return err
				}

				db, err := o.StateHandlers(o.StateNodes()).Component.GetPostgresDatabase(
					cfn.NewStackNamer().RDSPostgres(opts.ApplicationName, o.Declaration.Metadata.Name),
				)
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
		),
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

			app := fmt.Sprintf("%s-pgbouncer-%s", opts.ApplicationName, o.UserState.User.Username)

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
				o.Logger.Info("removing pgbouncer security group policy")

				perr := policyClient.Delete()
				if perr != nil {
					o.Logger.Warnf("deleting pgbouncer security group policy: %s", perr)

					err = perr
				}
			}()

			client := pgbouncer.New(&pgbouncer.Config{
				Name:                  app,
				Database:              opts.DatabaseName,
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
				Logger:                o.Logger,
			})

			err = client.Create()
			if err != nil {
				return err
			}

			defer func() {
				o.Logger.Info("removing pgbouncer pod")
				cerr := client.Delete()
				if cerr != nil {
					o.Logger.Warnf("deleting pgbouncer pod: %s", cerr)
					err = cerr
				}
			}()

			return err
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			metrics.Publish(metrics.Event{
				Category: metrics.CategoryCommandExecution,
				Action:   metrics.ActionForwardPostgres,
				Label:    metrics.LabelEnd,
			})

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.ApplicationName,
		"name",
		"n",
		"",
		"The name of the database instance to forward to",
	)

	flags.StringVarP(&opts.DatabaseName,
		"database",
		"d",
		"",
		"The name of the internal database name to use",
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
