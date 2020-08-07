package main

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/config"
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

// nolint: funlen
func buildShowCredentialsCommand(o *okctl.Okctl) *cobra.Command {
	opts := struct {
		Environment string
	}{}

	cmd := &cobra.Command{
		Use:   "credentials [env]",
		Short: "Show the location of the credentials",
		Long:  `This makes it possible to source the output from this command to run with kubectl`,
		Args:  cobra.ExactArgs(showCredentialsArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Environment = args[0]

			awsAccountID, err := o.AWSAccountID(opts.Environment)
			if err != nil {
				return err
			}

			return o.Initialise(opts.Environment, awsAccountID)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			appDir, err := o.GetAppDataDir()
			if err != nil {
				return err
			}

			kubeConfig := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterKubeConfig)
			awsConfig := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterAwsConfig)
			awsCredentials := path.Join(appDir, config.DefaultCredentialsDirName, o.ClusterName(opts.Environment), config.DefaultClusterAwsCredentials)

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
