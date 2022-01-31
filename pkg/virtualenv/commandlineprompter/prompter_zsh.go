package commandlineprompter

import (
	"strings"
)

type zshPrompter struct {
	clusterName string
	osEnvVars   map[string]string
}

// CreatePrompt returns environment variables that when set in zsh will show a command prompt.
// The warning is set in case something prevented the prompt to be set the expected way.
func (p *zshPrompter) CreatePrompt() (CommandLinePrompt, error) {
	ps1, overridePs1 := p.osEnvVars["OKCTL_PS1"]
	if !overridePs1 {
		ps1 = `%F{red}%~ %f%F{blue}($(venv_ps1 %env)%f) $ `
	}

	p.osEnvVars["PS1"] = strings.ReplaceAll(ps1, "%env", p.clusterName)

	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}
