package main

import (
	"fmt"
	"io"
	osPkg "os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/oslokommune/okctl/cmd/okctl/preruns"
	"github.com/oslokommune/okctl/pkg/metrics"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/kubeconfig"
	"github.com/oslokommune/okctl/pkg/virtualenv/shellgetter"

	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/storage"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/virtualenv"
	"github.com/spf13/cobra"
)

const (
	venvArgs = 0
)

func buildVenvCommand(o *okctl.Okctl) *cobra.Command { //nolint: funlen
	okctlEnvironment := commands.OkctlEnvironment{}

	cmd := &cobra.Command{
		Use:   "venv",
		Short: VenvShortDescription,
		Long:  VenvLongDescription,
		Args:  cobra.ExactArgs(venvArgs),
		PreRunE: preruns.PreRunECombinator(
			preruns.LoadUserData(o),
			preruns.InitializeMetrics(o),
			func(_ *cobra.Command, args []string) error {
				metrics.Publish(generateStartEvent(metrics.ActionVenv))

				e, err := venvPreRunE(o)
				if err != nil {
					return err
				}

				okctlEnvironment = e

				err = commands.ValidateBinaryVersionNotLessThanClusterVersion(o)
				if err != nil {
					return err
				}

				return nil
			},
		),
		RunE: func(_ *cobra.Command, args []string) error {
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

			err = o.FileSystem.MkdirAll(path.Dir(kubeConfigFile), 0o700)
			if err != nil {
				return err
			}

			err = o.FileSystem.WriteFile(kubeConfigFile, data, 0o640)
			if err != nil {
				return err
			}

			return venvRunE(o, okctlEnvironment)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			metrics.Publish(generateEndEvent(metrics.ActionVenv))

			return nil
		},
	}

	return cmd
}

func venvPreRunE(o *okctl.Okctl) (commands.OkctlEnvironment, error) {
	err := o.Initialise()
	if err != nil {
		return commands.OkctlEnvironment{}, err
	}

	okctlEnvironment, err := commands.GetOkctlEnvironment(o, declarationPath)
	if err != nil {
		return commands.OkctlEnvironment{}, err
	}

	return okctlEnvironment, nil
}

func venvRunE(o *okctl.Okctl, okctlEnvironment commands.OkctlEnvironment) error {
	venvOpts, tmpStorage, err := createVenvOpts(o.Host(), okctlEnvironment)
	if err != nil {
		return err
	}

	venv, err := virtualenv.CreateVirtualEnvironment(venvOpts)
	if err != nil {
		return fmt.Errorf("could not create virtual environment: %w", err)
	}

	defer func() {
		err = tmpStorage.Clean()
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = printWelcomeMessage(o.Out, venv, okctlEnvironment)
	if err != nil {
		return fmt.Errorf("could not print welcome message: %w", err)
	}

	shell := exec.Command(venv.ShellCommand) //nolint:gosec
	shell.Env = venv.Environ()
	shell.Stdout = o.Out
	shell.Stdin = o.In
	shell.Stderr = o.Err

	err = shell.Run()
	if err != nil {
		return fmt.Errorf("could not run shell: %w", err)
	}

	_, err = fmt.Fprintln(o.Out, "Exiting okctl virtual environment")
	if err != nil {
		return fmt.Errorf("could not print message: %w", err)
	}

	return nil
}

func createVenvOpts(host state.Host, okctlEnvironment commands.OkctlEnvironment) (commandlineprompter.CommandLinePromptOpts, *storage.TemporaryStorage, error) {
	currentUser, err := user.Current()
	if err != nil {
		return commandlineprompter.CommandLinePromptOpts{}, nil, fmt.Errorf("could not get current user: %w", err)
	}

	envVars := commands.GetVenvEnvVars(okctlEnvironment)

	homeDir, err := osPkg.UserHomeDir()
	if err != nil {
		return commandlineprompter.CommandLinePromptOpts{}, nil, fmt.Errorf("could not get user's home directory: %w", err)
	}

	venvOpts := commandlineprompter.CommandLinePromptOpts{
		Os:                   getOs(host),
		MacOsUserShellGetter: shellgetter.NewMacOsCmdGetter(homeDir),
		OsEnvVars:            envVars,
		EtcStorage:           storage.NewFileSystemStorage("/etc"),
		UserDirStorage:       storage.NewFileSystemStorage(okctlEnvironment.UserDataDir),
		UserHomeDirStorage:   storage.NewFileSystemStorage(homeDir),
		TmpStorage:           nil,
		ClusterName:          okctlEnvironment.ClusterName,
		CurrentUsername:      currentUser.Username,
	}

	tmpStorage, err := storage.NewTemporaryStorage()
	if err != nil {
		return commandlineprompter.CommandLinePromptOpts{}, nil, err
	}

	venvOpts.TmpStorage = tmpStorage

	return venvOpts, tmpStorage, nil
}

func getOs(host state.Host) shellgetter.Os {
	var os shellgetter.Os

	switch host.Os {
	case state.OsLinux:
		os = shellgetter.OsLinux
	case state.OsDarwin:
		os = shellgetter.OsDarwin
	default:
		os = shellgetter.OsUnknown
	}

	return os
}

type venvWelcomeMessage struct {
	ClusterName             string
	Environment             string
	KubectlPath             string
	AwsIamAuthenticatorPath string
	CommandPrompt           string
	VenvShellCmd            string
	VenvCommand             string
	Warning                 string
}

func printWelcomeMessage(stdout io.Writer, venv *virtualenv.VirtualEnvironment, opts commands.OkctlEnvironment) error {
	environ := venv.Environ()

	whichKubectl, err := getWhich("kubectl", environ)
	if err != nil {
		return fmt.Errorf("could not get executable 'kubectl': %w", err)
	}

	whichAwsIamAuthenticator, err := getWhich("aws-iam-authenticator", environ)
	if err != nil {
		return fmt.Errorf("could not get executable 'aws-iam-authenticator': %w", err)
	}

	params := venvWelcomeMessage{
		ClusterName:             opts.ClusterName,
		KubectlPath:             whichKubectl,
		AwsIamAuthenticatorPath: whichAwsIamAuthenticator,
		CommandPrompt:           "<directory> <okctl environment:kubernetes namespace>",
		VenvShellCmd:            aurora.Green(venv.ShellCommand).String(),
		VenvCommand:             aurora.Green("okctl venv").String(),
		Warning:                 venv.Warning,
	}
	template := `----------------- OKCTL -----------------
Cluster: {{ .ClusterName }}
Using kubectl: {{ .KubectlPath }}
Using aws-iam-authenticator: {{ .AwsIamAuthenticatorPath }}

Your shell is {{ .VenvShellCmd }} (override by setting environment variable OKCTL_SHELL), and your command prompt now shows
{{ .CommandPrompt }}

You can override the command prompt by setting the environment variable OKCTL_PS1 before running {{ .VenvCommand }}.
Or, if you do not want {{ .VenvCommand }} to modify your command prompt at all, set the environment variable
OKCTL_NO_PS1=true
`

	if len(venv.Warning) > 0 {
		template += "\n{{ .Warning }}\n"
	}

	template += "-----------------------------------------\n"

	msg, err := commands.GoTemplateToString(template, params)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(stdout, msg)
	if err != nil {
		return err
	}

	return nil
}

func getWhich(cmd string, env []string) (string, error) {
	c := exec.Command("which", cmd) //nolint:gosec
	c.Env = env

	o, err := c.Output()
	if err != nil {
		return "", err
	}

	which := string(o)
	which = strings.TrimSpace(which)

	return which, nil
}
