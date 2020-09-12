package main

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/kubeconfig"
	"sigs.k8s.io/yaml"

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

const showMsg = `
Now you can use %s to list nodes, pods, etc. Try out some commands:

$ %s get pods --all-namespaces
$ %s get nodes

This also requires %s, which you can add to your PATH from here:

%s

Optionally, install kubectl and aws-iam-authenticator to your
system from:

- https://kubernetes.io/docs/tasks/tools/install-kubectl/
- https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html

The installed version of kubectl needs to be within 2 versions of the
kubernetes cluster version, which is: %s.

We have also setup %s for continuous deployment, you can access
the UI at this URL by logging in with Github:

%s

`

// nolint: funlen gocognit
func buildShowCredentialsCommand(o *okctl.Okctl) *cobra.Command {
	opts := ShowCredentialsOpts{}

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

			meta := o.RepoStateWithEnv.GetMetadata()
			cluster := o.RepoStateWithEnv.GetCluster()

			opts.Repository = meta.Name
			opts.Region = meta.Region
			opts.AWSAccountID = cluster.AWSAccountID
			opts.Environment = cluster.Environment
			opts.ClusterName = cluster.Name

			err = opts.Validate()
			if err != nil {
				return errors.E(err, "failed to validate show credentials options")
			}

			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			outputDir, err := o.GetRepoOutputDir(opts.Environment)
			if err != nil {
				return err
			}

			appDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			kubeConfig := path.Join(appDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterKubeConfig)
			awsConfig := path.Join(appDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsConfig)
			awsCredentials := path.Join(appDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsCredentials)

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

			_, err = fmt.Fprintf(o.Out,
				"export AWS_CONFIG_FILE=%s\nexport AWS_SHARED_CREDENTIALS_FILE=%s\nexport AWS_PROFILE=default\nexport KUBECONFIG=%s\n",
				awsConfig,
				awsCredentials,
				kubeConfig,
			)
			if err != nil {
				return err
			}

			k, err := o.BinariesProvider.Kubectl(kubectl.Version)
			if err != nil {
				return err
			}

			a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(o.Err,
				showMsg,
				aurora.Green("kubectl"),
				k.BinaryPath,
				k.BinaryPath,
				aurora.Green("aws-iam-authenticator"),
				a.BinaryPath,
				aurora.Green("1.17"),
				aurora.Green("ArgoCD"),
				o.RepoStateWithEnv.GetArgoCD().SiteURL,
			)
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
