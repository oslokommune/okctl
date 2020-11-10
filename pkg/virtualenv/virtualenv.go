// Package virtualenv helps finding the environment variables needed to use a okctl cluster.
package virtualenv

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
	"github.com/oslokommune/okctl/pkg/virtualenv/shellgetter"
	"sort"
)

// VirtualEnvironment contains environment variables in a virtual environment.
type VirtualEnvironment struct {
	env          map[string]string
	Warning      string
	ShellCommand string
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

// Getenv returns the environment variable with the given key, and a bool indicating if the key was found or not.
func (v *VirtualEnvironment) Getenv(key string) (string, bool) {
	val, hasKey := v.env[key]
	return val, hasKey
}

// New returns a new virtual environment
func Create(opts commandlineprompter.CommandLinePromptOpts) (*VirtualEnvironment, error) {
	lsg := shellgetter.NewShellGetter(opts.OsEnvVars, opts.EtcStorage, opts.CurrentUsername)
	shell, err := lsg.Get()
	if err != nil {
		return nil, fmt.Errorf("could not get shell command: %w", err)
	}

	prompter, err := commandlineprompter.NewCommandLinePrompter(opts, shell.ShellType)
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
	}, nil
}
