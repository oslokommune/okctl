package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/cmd"
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
			shell := exec.Command(getShell(os.LookupEnv)) //nolint:gosec

			venv, err := virtualenv.GetVirtualEnvironment(&opts, os.Environ())
			if err != nil {
				return err
			}

			err = printExecutables(venv)
			if err != nil {
				return err
			}

			shell.Env = venv
			shell.Stdout = os.Stdout
			shell.Stdin = os.Stdin
			shell.Stderr = os.Stderr

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

func getShell(osLookupEnv func(key string) (string, bool)) string {
	shell, ok := osLookupEnv("OKCTL_SHELL")
	if ok {
		return shell
	}

	shell, ok = osLookupEnv("SHELL")
	if ok {
		return shell
	}

	return "/bin/sh"
}

func printExecutables(env []string) error {
	whichKubectl, err := getWhich("kubectl", env)
	if err != nil {
		return err
	}

	whichAwsIamAuthenticator, err := getWhich("aws-iam-authenticator", env)
	if err != nil {
		return err
	}

	fmt.Printf("Using kubectl: %s", whichKubectl)
	fmt.Printf("Using aws-iam-authenticator: %s", whichAwsIamAuthenticator)
	fmt.Println()

	return nil
}

func getWhich(cmd string, env []string) ([]byte, error) {
	c := exec.Command("which", cmd) //nolint:gosec
	c.Env = env

	return c.Output()
}
