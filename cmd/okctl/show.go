package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/virtualenv"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/kubeconfig"
	"sigs.k8s.io/yaml"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type showMessageOpts struct {
	VenvCmd                 string
	KubectlCmd              string
	KubectlPath             string
	AwsIamAuthenticatorCmd  string
	AwsIamAuthenticatorPath string
	K8sClusterVersion       string
	ArgoCD                  string
	ArgoCDURL               string
}

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

const showMsg = `
Tip: Run {{ .VenvCmd }} to run a shell with these environment variables set. Then you
can avoid using full paths to executables and modifying your PATH.

Now you can use {{ .KubectlCmd }} to list nodes, pods, etc. Try out some commands:

$ {{ .KubectlPath }} get pods --all-namespaces
$ {{ .KubectlPath }} get nodes

This also requires {{ .AwsIamAuthenticatorCmd }}, which you can add to your PATH from here:

{{ .AwsIamAuthenticatorPath }}

Optionally, install kubectl and aws-iam-authenticator to your
system from:

- https://kubernetes.io/docs/tasks/tools/install-kubectl/
- https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html

The installed version of kubectl needs to be within 2 versions of the
kubernetes cluster version, which is: {{ .K8sClusterVersion }}.

We have also setup {{ .ArgoCD }} for continuous deployment, you can access
the UI at this URL by logging in with Github:

{{ .ArgoCDURL }}

`

// nolint: funlen gocognit
func buildShowCredentialsCommand(o *okctl.Okctl) *cobra.Command {
	opts := virtualenv.VirtualEnvironmentOpts{}

	cmd := &cobra.Command{
		Use:   "credentials [env]",
		Short: "Show the location of the credentials",
		Long:  `This makes it possible to source the output from this command to run with kubectl`,
		Args:  cobra.ExactArgs(showCredentialsArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]

			err := o.InitialiseWithOnlyEnv(environment)
			if err != nil {
				return err
			}

			opts, err = virtualenv.GetVirtualEnvironmentOpts(o)

			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			osEnv := []string{"PATH=" + os.Getenv("PATH")}
			venv, err := virtualenv.GetVirtualEnvironment(&opts, osEnv)
			if err != nil {
				return err
			}

			for _, v := range venv {
				fmt.Fprintf(o.Out, "export %s\n", v)
			}

			outputDir, err := o.GetRepoOutputDir(opts.Environment)
			if err != nil {
				return err
			}

			appDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			kubeConfig := path.Join(appDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterKubeConfig)

			k, err := o.BinariesProvider.Kubectl(kubectl.Version)
			if err != nil {
				return err
			}

			a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
			if err != nil {
				return err
			}

			msg := showMessageOpts{
				VenvCmd:                 aurora.Green("okctl venv").String(),
				KubectlCmd:              aurora.Green("kubectl").String(),
				AwsIamAuthenticatorCmd:  aurora.Green("aws-iam-authenticator").String(),
				KubectlPath:             k.BinaryPath,
				AwsIamAuthenticatorPath: a.BinaryPath,
				K8sClusterVersion:       aurora.Green("1.17").String(),
				ArgoCD:                  aurora.Green("ArgoCD").String(),
				ArgoCDURL:               o.RepoStateWithEnv.GetArgoCD().SiteURL,
			}
			txt, err := goTemplateToString(showMsg, msg)
			if err != nil {
				return err
			}

			_, err = fmt.Print(o.Err, txt)
			if err != nil {
				return err
			}

			data, err := o.FileSystem.ReadFile(path.Join(outputDir, config.DefaultClusterBaseDir, config.DefaultClusterConfig))
			if err != nil {
				return err
			}

			clusterConfig := &v1alpha1.ClusterConfig{}

			err = yaml.Unmarshal(data, clusterConfig)
			if err != nil {
				return err
			}

			cfg, err := kubeconfig.New(clusterConfig, o.CloudProvider).Get()
			if err != nil {
				return err
			}

			data, err = cfg.Bytes()
			if err != nil {
				return err
			}

			err = o.FileSystem.WriteFile(kubeConfig, data, 0o644)
			if err != nil {
				return err
			}

			return err
		},
	}

	return cmd
}

func goTemplateToString(templateString string, data interface{}) (string, error) {
	tmpl, err := template.New("t").Parse(templateString)
	if err != nil {
		return "", err
	}

	tmplBuffer := new(bytes.Buffer)
	err = tmpl.Execute(tmplBuffer, data)

	if err != nil {
		return "", err
	}

	return tmplBuffer.String(), nil
}
