package main

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
	"io"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/storage"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/virtualenv"
	"github.com/spf13/cobra"
)

const (
	venvArgs = 1
)

const venvLong = `Runs a sub shell with all needed environmental variables set.

The variables are the same as shown in "okctl show credentials". The shell command to run is retrieved from the first
environment variable that is set of the following: $OKCTL_SHELL, $SHELL. If none is set, "/bin/sh" is used.

So to override, you can run for instance:

export OKCTL_SHELL=/bin/bash
okctl venv myenv
`

// nolint: gocyclo, funlen, gocognit
func buildVenvCommand(o *okctl.Okctl) *cobra.Command {
	credentialsOpts := commands.CredentialsOpts{}

	cmd := &cobra.Command{
		Use:   "venv ENV",
		Short: "Runs a virtual environment",
		Long:  venvLong,
		Args:  cobra.ExactArgs(venvArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			environment := args[0]

			err := validation.Validate(
				&environment,
				validation.Required,
				validation.Match(regexp.MustCompile("^[a-zA-Z]{3,64}$")).Error("the environment must consist of 3-64 characters (a-z, A-Z)"),
			)
			if err != nil {
				return err
			}

			err = o.InitialiseWithOnlyEnv(environment)
			if err != nil {
				return err
			}

			credentialsOpts, err = commands.GetCredentialsOpts(o)

			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			currentUser, err := user.Current()
			if err != nil {
				return fmt.Errorf("could not get current user: %w", err)
			}

			okctlEnvVars := commands.GetOkctlEnvVars(credentialsOpts)
			envVars := commands.MergeEnvVars(os.Environ(), okctlEnvVars)

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("could not get user's home directory: %w", err)
			}

			venvOpts := commandlineprompter.CommandLinePromptOpts{
				OsEnvVars:          envVars,
				EtcStorage:         storage.NewFileSystemStorage("/etc"),
				UserDirStorage:     storage.NewFileSystemStorage(credentialsOpts.UserDataDir),
				UserHomeDirStorage: storage.NewFileSystemStorage(homeDir),
				TmpStorage:         nil,
				Environment:        credentialsOpts.Environment,
				CurrentUsername:    currentUser.Username,
			}

			tmpStorage, err := storage.NewTemporaryStorage()
			if err != nil {
				return err
			}

			venvOpts.TmpStorage = tmpStorage

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

			err = printWelcomeMessage(o.Out, venv, credentialsOpts)
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
		},
	}

	return cmd
}

type venvWelcomeMessage struct {
	Environment             string
	KubectlPath             string
	AwsIamAuthenticatorPath string
	CommandPrompt           string
	VenvCommand             string
	Warning                 string
}

func printWelcomeMessage(stdout io.Writer, venv *virtualenv.VirtualEnvironment, opts commands.CredentialsOpts) error {
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
		Environment:             opts.Environment,
		KubectlPath:             whichKubectl,
		AwsIamAuthenticatorPath: whichAwsIamAuthenticator,
		CommandPrompt:           "<directory> <okctl environment:kubernetes namespace>",
		VenvCommand:             aurora.Green("okctl venv").String(),
		Warning:                 venv.Warning,
	}
	template := `----------------- OKCTL -----------------
Environment: {{ .Environment }}
Using kubectl: {{ .KubectlPath }}
Using aws-iam-authenticator: {{ .AwsIamAuthenticatorPath }}

Your command prompt now shows
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
