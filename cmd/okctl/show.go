package main

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/oslokommune/okctl/pkg/kubeconfig"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	showCredentialsArgs = 0
)

func buildShowCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: ShowCommandsShortDescription,
	}

	cmd.AddCommand(buildShowCredentialsCommand(o))

	return cmd
}

//nolint:funlen,gocognit
func buildShowCredentialsCommand(o *okctl.Okctl) *cobra.Command {
	okctlEnvironment := commands.OkctlEnvironment{}

	cmd := &cobra.Command{
		Use:   "credentials",
		Short: ShowShortDescription,
		Long:  ShowLongDescription,
		Args:  cobra.ExactArgs(showCredentialsArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			err := o.Initialise()
			if err != nil {
				return err
			}

			okctlEnvironment, err = commands.GetOkctlEnvironment(o, declarationPath)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			okctlEnvVars := commands.GetOkctlEnvVars(okctlEnvironment)

			for k, v := range okctlEnvVars {
				_, err := fmt.Fprintf(o.Out, "export %s=%s\n", k, v)
				if err != nil {
					return err
				}
			}

			k, err := o.BinariesProvider.Kubectl(kubectl.Version)
			if err != nil {
				return err
			}

			a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
			if err != nil {
				return err
			}

			argo, err := o.StateHandlers(o.StateNodes()).ArgoCD.GetArgoCD()
			if err != nil {
				return err
			}

			msg := commands.ShowMessageOpts{
				VenvCmd:                 aurora.Green("okctl venv").String(),
				KubectlCmd:              aurora.Green("kubectl").String(),
				AwsIamAuthenticatorCmd:  aurora.Green("aws-iam-authenticator").String(),
				KubectlPath:             k.BinaryPath,
				AwsIamAuthenticatorPath: a.BinaryPath,
				K8sClusterVersion:       aurora.Green("1.17").String(),
				ArgoCD:                  aurora.Green("ArgoCD").String(),
				ArgoCDURL:               argo.ArgoURL,
			}
			txt, err := commands.GoTemplateToString(commands.ShowMsg, msg)
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(o.Err, txt)
			if err != nil {
				return err
			}

			handlers := o.StateHandlers(o.StateNodes())

			cluster, err := handlers.Cluster.GetCluster(o.Declaration.Metadata.Name)
			if err != nil {
				return err
			}

			cfg, err := kubeconfig.New(cluster.Config, o.CloudProvider).Get()
			if err != nil {
				return fmt.Errorf("creating kubconfig: %w", err)
			}

			data, err := cfg.Bytes()
			if err != nil {
				return err
			}

			appDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			kubeConfigFile := path.Join(appDir, constant.DefaultCredentialsDirName, okctlEnvironment.ClusterName, constant.DefaultClusterKubeConfig)

			err = o.FileSystem.WriteFile(kubeConfigFile, data, 0o640)
			if err != nil {
				return err
			}

			return err
		},
	}

	return cmd
}
