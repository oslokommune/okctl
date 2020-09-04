package main

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/mishudark/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	showCredentialsArgs = 1
)

func buildShowCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show commands",
	}

	cmd.AddCommand(buildShowCredentialsCommand(o))

	return cmd
}

// ShowCredentialsOpts contains the required inputs
type ShowCredentialsOpts struct {
	Region       string
	AWSAccountID string
	Environment  string
	Repository   string
	ClusterName  string
}

// Validate the inputs
func (o *ShowCredentialsOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: funlen
func buildShowCredentialsCommand(o *okctl.Okctl) *cobra.Command {
	opts := ShowCredentialsOpts{}

	cmd := &cobra.Command{
		Use:   "credentials [env]",
		Short: "Show the location of the credentials",
		Long:  `This makes it possible to source the output from this command to run with kubectl`,
		Args:  cobra.ExactArgs(showCredentialsArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]

			opts.Region = o.Region()
			opts.AWSAccountID = o.AWSAccountID(environment)
			opts.Environment = environment
			opts.Repository = o.RepoData.Name
			opts.ClusterName = o.ClusterName(environment)

			err := opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate show credentials options")
			}

			return o.Initialise(opts.Environment, opts.AWSAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			appDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			kubeConfig := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterKubeConfig)
			awsConfig := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterAwsConfig)
			awsCredentials := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterAwsCredentials)

			h := &helm.Config{
				HelmPluginsDirectory: path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
				HelmRegistryConfig:   path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
				HelmRepositoryConfig: path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
				HelmRepositoryCache:  path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
				HelmBaseDir:          path.Join(appDir, config.DefaultHelmBaseDir),
				Debug:                o.Debug,
			}

			for k, v := range h.Envs() {
				if k == "HOME" || k == "PATH" {
					continue
				}

				_, err = io.Copy(o.Out, strings.NewReader(fmt.Sprintf("export %s=%s\n", k, v)))
				if err != nil {
					return err
				}
			}

			_, err = io.Copy(
				o.Out,
				strings.NewReader(
					fmt.Sprintf(
						"export AWS_CONFIG_FILE=%s\nexport AWS_SHARED_CREDENTIALS_FILE=%s\nexport AWS_PROFILE=default\nexport KUBECONFIG=%s\n",
						awsConfig,
						awsCredentials,
						kubeConfig,
					),
				),
			)

			return err
		},
	}

	return cmd
}
