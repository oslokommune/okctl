package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/commands"
	cmd "github.com/oslokommune/okctl/pkg/commands"

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

// nolint: funlen
func buildVenvCommand(o *okctl.Okctl) *cobra.Command {
	opts := virtualenv.VirtualEnvironmentOpts{}

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

			opts, err = virtualenv.GetVirtualEnvironmentOpts(o)

			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			shell := exec.Command(cmd.GetShell(os.LookupEnv)) //nolint:gosec

			venv, err := virtualenv.GetVirtualEnvironment(&opts, os.Environ())
			if err != nil {
				return err
			}

			err = printWelcomeMessage(o, venv, &opts)
			if err != nil {
				return err
			}

			shell.Env = venv
			shell.Stdout = o.Out
			shell.Stdin = o.In
			shell.Stderr = o.Err

			err = shell.Run()
			if err != nil {
				log.Fatalf("Command failed: %v", err)
			}

			fmt.Println("Exiting okctl virtual environment")

			return nil
		},
	}

	return cmd
}

type VenvWelcomeMessage struct {
	Environment             string
	KubectlPath             string
	AwsIamAuthenticatorPath string
	CommandPrompt           string
	VenvCommand             string
}

func printWelcomeMessage(o *okctl.Okctl, env []string, opts *virtualenv.VirtualEnvironmentOpts) error {
	whichKubectl, err := getWhich("kubectl", env)
	if err != nil {
		return err
	}

	whichAwsIamAuthenticator, err := getWhich("aws-iam-authenticator", env)
	if err != nil {
		return err
	}

	params := VenvWelcomeMessage{
		Environment:             opts.Environment,
		KubectlPath:             whichKubectl,
		AwsIamAuthenticatorPath: whichAwsIamAuthenticator,
		CommandPrompt:           "<directory> <git branch> <okctl environment:kubernetes namespace>",
		VenvCommand:             aurora.Green("okctl venv").String(),
	}
	template := `----------------- OKCTL -----------------
Environment: {{ .Environment }}
Using kubectl: {{ .KubectlPath }}
Using aws-iam-authenticator: {{ .AwsIamAuthenticatorPath }}

Your command prompt now shows
{{ .CommandPrompt }}
You can override the command prompt by setting the environment variable OKCTL_PS1 before running {{ .VenvCommand }}

-----------------------------------------
`

	msg, err := commands.GoTemplateToString(template, params)
	if err != nil {
		return err
	}

	fmt.Fprint(o.Out, msg)

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
