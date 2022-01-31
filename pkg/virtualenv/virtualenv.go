// Package virtualenv helps finding the environment variables needed to use a okctl cluster.
package virtualenv

import (
	"fmt"
	"sort"

	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
)

// VirtualEnvironment contains environment variables in a virtual environment.
type VirtualEnvironment struct {
	env          map[string]string
	Warning      string
	ShellCommand string
	ShellArgs    []string
}

// Environ returns all environment variables in the virtual environment, on the form
// []string { "KEY1=VALUE1", "KEY2=VALUE2", ... }
// This is the same form as os.Environ.
func (v *VirtualEnvironment) Environ() []string {
	venvs := make([]string, 0, len(v.env))

	for k, v := range v.env {
		keyEqualsValue := fmt.Sprintf("%s=%s", k, v)
		venvs = append(venvs, keyEqualsValue)
	}

	sort.Strings(venvs)

	return venvs
}

// CreateVirtualEnvironment returns an environment for running a shell
func CreateVirtualEnvironment(opts commandlineprompter.CommandLinePromptOpts) (*VirtualEnvironment, error) {
	shellGetter := commandlineprompter.NewShellGetter(opts)

	shell, err := shellGetter.Get()
	if err != nil {
		return nil, fmt.Errorf("could not get shell: %w", err)
	}

	prompter, err := commandlineprompter.New(opts, shell.ShellType)
	if err != nil {
		return nil, fmt.Errorf("could not create command line prompter: %w", err)
	}

	commandLinePrompt, err := prompter.CreatePrompt()
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	return &VirtualEnvironment{
		env:          commandLinePrompt.Env,
		Warning:      commandLinePrompt.Warning,
		ShellCommand: shell.Command,
		ShellArgs:    shell.Args,
	}, nil
}
