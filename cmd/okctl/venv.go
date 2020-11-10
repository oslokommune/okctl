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
				return err
			}

			// 10 SHOW START
			// OG: Samme issue videre for å avgjøre hvilken virtual environment som skal lages
			// Sette PATH=$PATH:/path/kubectl:/path/awsIam
			// 10 SHOW END

			// AGH: bør kanskje lage shell getter her, fordi det er to mulige veier
			// alt1: hente shell fra OKCTL_SHELL
			// alt2: hente shell fra /etc/passwd

			// Lese OKCTL_NO_PS1. Hvis ikke satt:
			//
			// Lag ~/.okctl/venv_ps1 fil hvis ikke finnes
			// For begge: Sett PATH=$PATH:/path/to/venv_ps1_dir
			//
			// 1. for bash:
			//   Sette PROMPT_COMMAND = OKCTL_PS1 eller default ps1
			// 2. for zsh
			//   Skrive /tmp/.zshrc temp fil
			//	   Hvis ZDOTDIR finnes
			//	     [NEI] write: source ZDOTDIR/.zshrc - funker ikke, for da kan vi ikke sette ZDOTDIR selv
			//       VURDER: Skriv "WARNING: Could not set command prompt (PS1) because ZDOTDIR is already set.
			//       Either start okctl venv with no ZDOTDIR set, or set environment variable OKCTL_NO_PS1=true to get
			//       rid of this message.
			//       og drit å gjør noe (dvs, ikke sett ZDOTDIR, og ikke skriv til /tmp/.zshrc heller)
			//
			//     Else if ~/.zshrc finnes
			//       write: source ~/.zshrc
			//     ++, se cmdprompt_zsh.go

			//   Sette ZDOTDIR = OKCTL_PS1 eller default ps1
			//

			okctlEnvVars := commands.GetOkctlEnvVars(credentialsOpts)
			envVars := commands.MergeEnvVars(os.Environ(), okctlEnvVars)

			venvOpts := commandlineprompter.CommandLinePromptOpts{
				OsEnvVars:       envVars,
				EtcStorage:      storage.NewFileSystemStorage("/etc"),
				UserDirStorage:  storage.NewFileSystemStorage(credentialsOpts.UserDataDir),
				TmpStorage:      nil,
				Environment:     credentialsOpts.Environment,
				CurrentUsername: currentUser.Username,
			}

			tmpStorage, err := storage.NewTemporaryStorage()
			if err != nil {
				return err
			}

			venvOpts.TmpStorage = tmpStorage

			venv, err := virtualenv.Create(venvOpts)
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

	whichKubectl, err := getWhichWithEnv("kubectl", environ)
	if err != nil {
		return fmt.Errorf("could not get executable 'kubectl': %w", err)
	}

	whichAwsIamAuthenticator, err := getWhichWithEnv("aws-iam-authenticator", environ)
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

func getWhichWithEnv(cmd string, env []string) (string, error) {
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
