package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/logging"

	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
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

const defaultPostgreSQLPort = 5432

type forwardPostgresOpts struct {
	ID                     api.ID
	ClusterDeclarationPath string
	ApplicationName        string
	DatabaseName           string
	Namespace              string
	ConfigMapName          string
	SecretName             string
	SecurityGroup          string
	Username               string
	PasswordFile           string
	ListenPort             int32
}

// Validate the inputs
func (o *forwardPostgresOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ClusterDeclarationPath, validation.Required),
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

	initEnvOpts := hooks.InitializeEnvironmentOpts{
		InitSignalHandling: false,
	}

	cmd := &cobra.Command{
		Use:   "postgres",
		Short: ForwardPostgresShortDescription,
		Long:  ForwardPostgresLongDescription,
		Args:  cobra.ExactArgs(0), // nolint: gomnd
		// We override PersistentPreRunE in order to override the inherited signal handling, because we want to clean up pods if
		// user press CTRL+C. Without this, pods are left behind running.
		PersistentPreRunE: hooks.InitializeEnvironmentWithOpts(o, initEnvOpts),
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionForwardPostgres),
			hooks.LoadClusterDeclaration(o, &opts.ClusterDeclarationPath),
			hooks.InitializeOkctl(o),
			hooks.DownloadState(o, false),
			hooks.VerifyClusterExistsInState(o),
			hooks.ValidatePostgresApplicationName(o, &opts.ApplicationName),
			func(_ *cobra.Command, _ []string) error {
				err := commands.ValidateBinaryVersionNotLessThanClusterVersion(o)
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
			services, err := o.ClientServices(o.StateHandlers(o.StateNodes()))
			if err != nil {
				return fmt.Errorf("acquiring client services: %w", err)
			}

			clusterSecurityGroupID, err := services.Cluster.GetClusterSecurityGroupID(
				context.Background(),
				client.GetClusterSecurityGroupIDOpts{ID: opts.ID},
			)
			if err != nil {
				return fmt.Errorf("acquiring cluster security group ID: %w", err)
			}

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
					[]string{
						opts.SecurityGroup,
						clusterSecurityGroupID.Value,
					},
				),
				config,
			)

			_, err = policyClient.Create()
			if err != nil {
				return err
			}

			pgBouncerClient := pgbouncer.New(&pgbouncer.Config{
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

			activateSigintCatch(policyClient.Delete, pgBouncerClient.Delete)

			err = pgBouncerClient.Create()
			if err != nil {
				return err
			}

			return err
		},
		PostRunE: hooks.RunECombinator(
			hooks.ClearLocalState(o),
			hooks.EmitEndCommandExecutionEvent(metrics.ActionForwardPostgres),
		),
	}
	addAuthenticationFlags(cmd)
	addClusterDeclarationPathFlag(cmd, &opts.ClusterDeclarationPath)

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
		defaultPostgreSQLPort,
		"The port on your local machine traffic will be forwarded to",
	)

	return cmd
}

func activateSigintCatch(teardownFns ...func() error) {
	log := logging.GetLogger("forward postgres", "activateSigintCatch")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info("removing pgbouncer security group policy")

		for _, teardownFn := range teardownFns {
			err := teardownFn()
			if err != nil {
				log.Warn(err.Error())
			}
		}
	}()
}
