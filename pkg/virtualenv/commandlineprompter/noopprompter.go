package commandlineprompter

type noopPrompter struct {
	osEnvVars map[string]string
}

func (p *noopPrompter) CreatePrompt() (CommandLinePrompt, error) {
	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}
