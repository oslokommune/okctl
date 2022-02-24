package main

import (
	"fmt"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"
	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	showCredentialsArgs = 0
)

type showCredentialsOpts struct {
	ClusterDeclarationPath string
}

func buildShowCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: ShowCommandsShortDescription,
	}

	cmd.AddCommand(buildShowCredentialsCommand(o))
	addAuthenticationFlags(cmd)

	return cmd
}

//nolint:funlen,gocognit,gocyclo
func buildShowCredentialsCommand(o *okctl.Okctl) *cobra.Command {
	var (
		okctlEnvironment = commands.OkctlEnvironment{}
		err              error
	)

	opts := &showCredentialsOpts{}

	cmd := &cobra.Command{
		Use:   "credentials",
		Short: ShowShortDescription,
		Long:  ShowLongDescription,
		Args:  cobra.ExactArgs(showCredentialsArgs),
		PreRunE: hooks.RunECombinator(
			hooks.LoadUserData(o),
			hooks.InitializeMetrics(o),
			hooks.EmitStartCommandExecutionEvent(metrics.ActionShowCredentials),
			hooks.LoadClusterDeclaration(o, &opts.ClusterDeclarationPath),
			hooks.InitializeOkctl(o),
			hooks.DownloadState(o, false),
			hooks.VerifyClusterExistsInState(o),
			hooks.WriteKubeConfig(o),
			func(_ *cobra.Command, args []string) error {
				okctlEnvironment, err = commands.GetOkctlEnvironment(o, opts.ClusterDeclarationPath)
				if err != nil {
					return err
				}

				err := commands.ValidateBinaryVersionNotLessThanClusterVersion(o)
				if err != nil {
					return err
				}

				return nil
			},
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			okctlEnvVars, err := commands.GetOkctlEnvVars(okctlEnvironment)
			if err != nil {
				return err
			}

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

			return err
		},
		PostRunE: hooks.RunECombinator(
			hooks.ClearLocalState(o),
			hooks.EmitEndCommandExecutionEvent(metrics.ActionShowCredentials),
		),
	}
	addClusterDeclarationPathFlag(cmd, &opts.ClusterDeclarationPath)

	return cmd
}
