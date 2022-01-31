package commandlineprompter

type zshPrompter struct {
	clusterName string
	osEnvVars   map[string]string
}

// CreatePrompt returns environment variables that when set in zsh will show a command prompt.
// The warning is set in case something prevented the prompt to be set the expected way.
func (p *zshPrompter) CreatePrompt() (CommandLinePrompt, error) {
	p.osEnvVars["PS1"] = "%F{red}%~ %f%F{blue}($(venv_ps1 " + p.clusterName + ")%f) $ "

	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}
