package commandlineprompter

import (
	"fmt"
	"strings"
)

type bashPrompter struct {
	clusterName string
	osEnvVars   map[string]string
}

func (p *bashPrompter) CreatePrompt() (CommandLinePrompt, error) {
	ps1, overridePs1 := p.osEnvVars["OKCTL_PS1"]
	if overridePs1 {
		withEnv := strings.ReplaceAll(ps1, "%env", p.clusterName)
		p.osEnvVars["PROMPT_COMMAND"] = fmt.Sprintf(`PS1="%s"`, withEnv)
	} else {
		p.osEnvVars["PROMPT_COMMAND"] = fmt.Sprintf(`PS1="\[\e[0;31m\]\w \[\e[0;34m\](\$(venv_ps1 %s)) \[\e[0m\]\$ "`, p.clusterName)
	}

	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}
