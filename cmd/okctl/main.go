package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/oslokommune/okctl/cmd/okctl/auth"
	"github.com/oslokommune/okctl/cmd/okctl/hooks"

	"github.com/oslokommune/okctl/pkg/config/constant"

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

		err := hooks.GracefullyTearDownState(o)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "gracefully cleaning up state: %s", err.Error())
		}
	}()

	err := cmd.Execute()
	if err != nil {
		exitCode = 1
	}
}

func buildRootCommand() (*cobra.Command, *okctl.Okctl) {
	o := okctl.New()
	if err := o.InitLogging(); err != nil {
		fmt.Fprintln(os.Stderr, "Error configuring logging:", err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:               "okctl",
		Short:             OkctlShortDescription,
		Long:              OkctlLongDescription,
		SilenceUsage:      true,
		PersistentPreRunE: hooks.InitializeEnvironment(o),
	}

	addAvailableCommands(cmd, o)

	return cmd, o
}

func addAvailableCommands(cmd *cobra.Command, o *okctl.Okctl) {
	cmd.AddCommand(buildApplyCommand(o))
	cmd.AddCommand(buildCompletionCommand(o))
	cmd.AddCommand(buildDeleteCommand(o))
	cmd.AddCommand(buildScaffoldCommand(o))
	cmd.AddCommand(buildShowCommand(o))
	cmd.AddCommand(buildVenvCommand(o))
	cmd.AddCommand(buildForwardCommand(o))
	cmd.AddCommand(buildVersionCommand(o))
	cmd.AddCommand(buildUpgradeCommand(o))
	cmd.AddCommand(buildMaintenanceCommand(o))
}

// Add the common authentication flags used throughout the application.
// Each sub-command needs to apply this individual according to needs
func addAuthenticationFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP((*string)(&auth.AwsCredentialsType),
		"aws-credentials-type",
		"a",
		getWithDefault(os.Getenv, constant.EnvAWSCredentialsType, context.AWSCredentialsTypeSAML),
		fmt.Sprintf(
			"The form of authentication to use for AWS. Possible values: [%s]",
			strings.Join(auth.GetAwsCredentialsTypes(), ","),
		),
	)
	cmd.PersistentFlags().StringVarP((*string)(&auth.GithubCredentialsType),
		"github-credentials-type",
		"g",
		getWithDefault(os.Getenv, constant.EnvGithubCredentialsType, context.GithubCredentialsTypeDeviceAuthentication),
		fmt.Sprintf(
			"The form of authentication to use for Github. Possible values: [%s]",
			strings.Join(auth.GetGithubCredentialsTypes(), ","),
		),
	)
}

// Add a '-c/--cluster-declaration' flag for those that require this
// Each sub-command needs to apply this individual according to needs
func addClusterDeclarationPathFlag(cmd *cobra.Command, clusterDeclarationPath *string) {
	cmd.Flags().StringVarP(clusterDeclarationPath,
		"cluster-declaration",
		"c",
		os.Getenv(constant.EnvClusterDeclaration),
		usageApplyClusterFile,
	)
}

func getWithDefault(getter func(key string) string, key string, defaultValue string) string {
	rawValue := getter(key)

	if rawValue == "" {
		return defaultValue
	}

	return rawValue
}
