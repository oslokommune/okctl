package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"

	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadRepoData(o *okctl.Okctl, configFile string, _ *cobra.Command) error {
	o.RepoDataLoader = load.RepoDataFromConfigFile(configFile)

	return o.LoadRepoData()
}

func loadUserData(o *okctl.Okctl, cmd *cobra.Command) error {
	userDataNotFound := load.CreateOnUserDataNotFound()

	o.UserDataLoader = load.UserDataFromFlagsEnvConfigDefaults(cmd, userDataNotFound)

	return o.LoadUserData()
}

//nolint:funlen,govet
func buildRootCommand() *cobra.Command {
	var outputFormat, declarationPath string

	o := okctl.New()

	cmd := &cobra.Command{
		Use:   "okctl",
		Short: "Opinionated and effortless infrastructure and application management",
		Long: `A highly opinionated CLI for creating a Kubernetes cluster in AWS with
a set of applications that ensure tighter integration between AWS and
Kubernetes, e.g., aws-alb-ingress-controller, external-secrets, etc.

Also comes pre-configured with ArgoCD for managing deployments, etc.
We also use the prometheus-operator for ensuring metrics and logs are
being captured. Together with slack and slick.`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == cobra.ShellCompRequestCmd {
				return nil
			}

			enableServiceUserAuthentication(o)

			var err error

			if len(declarationPath) == 0 {
				return fmt.Errorf("declaration must be provided")
			}

			declarationPath, err = filepath.Abs(declarationPath)
			if err != nil {
				return fmt.Errorf("converting declaration path to absolute path: %w", err)
			}

			err = loadUserData(o, cmd)
			if err != nil {
				return fmt.Errorf("loading application data: %w", err)
			}

			err = loadRepoData(o, declarationPath, cmd)
			if err != nil {
				if errors.Is(err, git.ErrRepositoryNotExists) {
					return fmt.Errorf("okctl needs to be run inside a Git repository (okctl outputs " +
						"various configuration files that will be stored here)")
				}

				return fmt.Errorf("loading repository data: %w", err)
			}

			o.Out = cmd.OutOrStdout()
			o.Err = cmd.OutOrStderr()

			o.SetFormat(core.EncodeResponseType(outputFormat))

			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				os.Exit(1)
			}()

			return nil
		},
	}

	cmd.AddCommand(buildApplyCommand(o))
	cmd.AddCommand(buildCompletionCommand(o))
	cmd.AddCommand(buildCreateCommand(o))
	cmd.AddCommand(buildDeleteCommand(o))
	cmd.AddCommand(buildScaffoldCommand(o))
	cmd.AddCommand(buildShowCommand(o))
	cmd.AddCommand(buildVenvCommand(o))
	cmd.AddCommand(buildAttachCommand(o))
	cmd.AddCommand(buildForwardCommand(o))
	cmd.AddCommand(buildVersionCommand(o))

	f := cmd.Flags()
	f.StringVarP(&outputFormat, "output", "o", "text",
		"The format of the output returned to the user")

	cmd.PersistentFlags().StringVarP(&declarationPath,
		"cluster-declaration",
		"c",
		os.Getenv(fmt.Sprintf("%s_%s", constant.EnvPrefix, constant.EnvClusterDeclaration)),
		"The cluster declaration you want to use",
	)
	cmd.PersistentFlags().StringVarP(&awsCredentialsType,
		"aws-credentials-type",
		"a",
		context.AWSCredentialsTypeSAML,
		fmt.Sprintf(
			"The form of authentication to use for AWS. Possible values: [%s,%s]",
			context.AWSCredentialsTypeSAML,
			context.AWSCredentialsTypeAccessKey,
		),
	)
	cmd.PersistentFlags().StringVarP(&githubCredentialsType,
		"github-credentials-type",
		"g",
		context.GithubCredentialsTypeDeviceAuthentication,
		fmt.Sprintf(
			"The form of authentication to use for Github. Possible values: [%s,%s]",
			context.GithubCredentialsTypeDeviceAuthentication,
			context.GithubCredentialsTypeToken,
		),
	)

	return cmd
}
