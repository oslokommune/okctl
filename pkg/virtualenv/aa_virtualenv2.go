package virtualenv

import (
	"fmt"
)

// New returns a new virtual environment
func Create(opts VirtualEnvironmentOpts) (*VirtualEnvironment, error) {
	lsg := newShellGetter(opts.OsEnvVars, opts.EtcStorage, opts.CurrentUsername)
	shell, err := lsg.get()
	if err != nil {
		return nil, fmt.Errorf("could not get shell command: %w", err)
	}

	prompter, err := newCommandLinePrompter(opts, shell.shellType)
	if err != nil {
		return nil, fmt.Errorf("could not create command line prompter: %w", err)
	}

	commandLinePrompt, err := prompter.createPrompt()
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	return &VirtualEnvironment{
		env:          commandLinePrompt.Env,
		Warning:      commandLinePrompt.Warning,
		ShellCommand: shell.command,
	}, nil
}
