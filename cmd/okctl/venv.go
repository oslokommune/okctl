package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
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

			userDataDir, err := o.GetUserDataDir()
			if err != nil {
				return err
			}

			ps1Dir, err := createPs1ExecutableIfNotExists(userDataDir)
			if err != nil {
				return err
			}

			opts, err = virtualenv.GetVirtualEnvironmentOptsWithPs1(o, ps1Dir)

			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			shellCmd, err := virtualenv.GetShellCmd(os.LookupEnv, storage.NewFileSystemStorage("/etc"))
			if err != nil {
				return err
			}

			venv, err := virtualenv.GetVirtualEnvironment(&opts, os.Environ())
			if err != nil {
				return err
			}

			noPs1, noPs1Exists := os.LookupEnv("OKCTL_NO_PS1")
			if !noPs1Exists || (noPs1Exists && strings.ToLower(strings.TrimSpace(noPs1)) != "true") {
				switch {
				case virtualenv.ShellIsBash(shellCmd):
					virtualenv.SetCmdPromptBash(&opts, venv)
				case virtualenv.ShellIsZsh(shellCmd):
					zshrcTmpWriteStorage, err := setCmdPromptZsh(&opts, venv)
					if err != nil {
						return err
					}

					defer func() {
						err = zshrcTmpWriteStorage.Clean()
						if err != nil {
							fmt.Println(err)
						}
					}()
				default:
					// We don't support any other shells for now.
				}
			}

			err = printWelcomeMessage(o.Out, venv.Environ(), &opts)
			if err != nil {
				return err
			}

			shell := exec.Command(shellCmd) //nolint:gosec
			shell.Env = venv.Environ()
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

func setCmdPromptZsh(opts *virtualenv.VirtualEnvironmentOpts, venv *virtualenv.VirtualEnvironment) (*storage.TemporaryStorage, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}

	zshrcReadStorage := storage.NewFileSystemStorage(userHomeDir)

	zshrcTmpWriteStorage, err := storage.NewTemporaryStorage()
	if err != nil {
		return nil, fmt.Errorf("could not create temporary storage: %w", err)
	}

	err = virtualenv.SetCmdPromptZsh(opts, venv, zshrcReadStorage, zshrcTmpWriteStorage)
	if err != nil {
		return nil, fmt.Errorf("could not set command prompt for virtul environment: %w", err)
	}

	return zshrcTmpWriteStorage, nil
}

func createPs1ExecutableIfNotExists(userDataDir string) (string, error) {
	store := storage.NewFileSystemStorage(userDataDir)

	ps1Dirname, err := virtualenv.CreatePs1ExecutableIfNotExists(store)
	if err != nil {
		return "", err
	}

	return path.Join(userDataDir, ps1Dirname), nil
}

type venvWelcomeMessage struct {
	Environment             string
	KubectlPath             string
	AwsIamAuthenticatorPath string
	CommandPrompt           string
	VenvCommand             string
}

func printWelcomeMessage(stdout io.Writer, env []string, opts *virtualenv.VirtualEnvironmentOpts) error {
	whichKubectl, err := getWhichWithEnv("kubectl", env)
	if err != nil {
		return err
	}

	whichAwsIamAuthenticator, err := getWhichWithEnv("aws-iam-authenticator", env)
	if err != nil {
		return err
	}

	params := venvWelcomeMessage{
		Environment:             opts.Environment,
		KubectlPath:             whichKubectl,
		AwsIamAuthenticatorPath: whichAwsIamAuthenticator,
		CommandPrompt:           "<directory> <okctl environment:kubernetes namespace>",
		VenvCommand:             aurora.Green("okctl venv").String(),
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
-----------------------------------------
`

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
