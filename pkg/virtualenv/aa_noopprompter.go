package virtualenv

type noopPrompter struct {
	osEnvVars map[string]string
}

func (p *noopPrompter) createPrompt() (CommandLinePrompt, error) {
	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}
