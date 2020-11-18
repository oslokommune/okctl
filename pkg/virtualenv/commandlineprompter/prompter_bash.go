package commandlineprompter

import (
	"fmt"
	"strings"
)

type bashPrompter struct {
	environment string
	osEnvVars   map[string]string
}

func (p *bashPrompter) CreatePrompt() (CommandLinePrompt, error) {
	ps1, overridePs1 := p.osEnvVars["OKCTL_PS1"]
	if overridePs1 {
		withEnv := strings.ReplaceAll(ps1, "%env", p.environment)
		p.osEnvVars["PROMPT_COMMAND"] = fmt.Sprintf(`PS1="%s"`, withEnv)
	} else {
		p.osEnvVars["PROMPT_COMMAND"] = fmt.Sprintf(`PS1="\[\e[0;31m\]\w \[\e[0;34m\](\$(venv_ps1 %s)) \[\e[0m\]\$ "`, p.environment)
	}

	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}
