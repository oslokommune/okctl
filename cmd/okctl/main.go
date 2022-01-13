package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/oslokommune/okctl/cmd/okctl/hooks"

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
	cmd, o := buildRootCommand()
	exitCode := 0

	defer func() {
		os.Exit(exitCode)
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), r)

			exitCode = 1
		}

		err := gracefullyTearDownState(o)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "gracefully cleaning up state: %s", err.Error())
		}
	}()

	err := cmd.Execute()
	if err != nil {
		exitCode = 1
	}
}

func gracefullyTearDownState(o *okctl.Okctl) error {
	err := hooks.UploadState(o)(nil, nil)
	if err != nil {
		switch {
		case errors.Is(err, hooks.ErrNotInitialized):
			return nil
		case errors.Is(err, hooks.ErrImmutable):
			return nil
		case errors.Is(err, hooks.ErrNotFound):
			return nil
		default:
			return fmt.Errorf("uploading state: %w", err)
		}
	}

	err = hooks.ReleaseStateLock(o)(nil, nil)
	if err != nil {
		return fmt.Errorf("releasing state lock: %w", err)
	}

	err = hooks.ClearLocalState(o)(nil, nil)
	if err != nil {
		return fmt.Errorf("clearing local state: %w", err)
	}

	return nil
}

func loadRepoData(o *okctl.Okctl, declarationPath string, _ *cobra.Command) error {
	o.RepoDataLoader = load.RepoDataFromConfigFile(declarationPath)

	return o.LoadRepoData()
}

func loadUserData(o *okctl.Okctl, cmd *cobra.Command) error {
	userDataNotFound := load.CreateOnUserDataNotFound()

	o.UserDataLoader = load.UserDataFromFlagsEnvConfigDefaults(cmd, userDataNotFound)

	return o.LoadUserData()
}

var declarationPath string //nolint:gochecknoglobals

//nolint:funlen,govet
func buildRootCommand() (*cobra.Command, *okctl.Okctl) {
	var outputFormat string

	o := okctl.New()
	if err := o.InitLogging(); err != nil {
		fmt.Fprintln(os.Stderr, "Error configuring logging:", err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:          "okctl",
		Short:        OkctlShortDescription,
		Long:         OkctlLongDescription,
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

				err = gracefullyTearDownState(o)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "gracefully up state: %s", err.Error())
				}

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
	cmd.AddCommand(buildForwardCommand(o))
	cmd.AddCommand(buildVersionCommand(o))
	cmd.AddCommand(buildUpgradeCommand(o))
	cmd.AddCommand(buildMaintenanceCommand(o))

	f := cmd.Flags()
	f.StringVarP(&outputFormat, "output", "o", "text",
		"The format of the output returned to the user")

	cmd.PersistentFlags().StringVarP(&declarationPath,
		"cluster-declaration",
		"c",
		os.Getenv(constant.EnvClusterDeclaration),
		"The cluster declaration you want to use",
	)
	cmd.PersistentFlags().StringVarP(&awsCredentialsType,
		"aws-credentials-type",
		"a",
		getWithDefault(os.Getenv, constant.EnvAWSCredentialsType, context.AWSCredentialsTypeSAML),
		fmt.Sprintf(
			"The form of authentication to use for AWS. Possible values: [%s,%s]",
			context.AWSCredentialsTypeSAML,
			context.AWSCredentialsTypeAccessKey,
		),
	)
	cmd.PersistentFlags().StringVarP(&githubCredentialsType,
		"github-credentials-type",
		"g",
		getWithDefault(os.Getenv, constant.EnvGithubCredentialsType, context.GithubCredentialsTypeDeviceAuthentication),
		fmt.Sprintf(
			"The form of authentication to use for Github. Possible values: [%s,%s]",
			context.GithubCredentialsTypeDeviceAuthentication,
			context.GithubCredentialsTypeToken,
		),
	)

	return cmd, o
}

func getWithDefault(getter func(key string) string, key string, defaultValue string) string {
	rawValue := getter(key)

	if rawValue == "" {
		return defaultValue
	}

	return rawValue
}
