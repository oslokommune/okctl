package virtualenv

import (
	"fmt"
	"strings"
)

// ShellIsBash returns true if provided command will run bash
func ShellIsBash(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "bash")
}

// SetCmdPromptBash makes bash show a custom command prompt. It does so by setting the environment variable
// PROMPT_COMMAND.
func SetCmdPromptBash(opts *VirtualEnvironmentOpts, venv *VirtualEnvironment) {
	ps1, overridePs1 := venv.env["OKCTL_PS1"]
	if overridePs1 {
		venv.env["PROMPT_COMMAND"] = fmt.Sprintf("PS1=%s", ps1)
	} else {
		venv.env["PROMPT_COMMAND"] = fmt.Sprintf(`PS1="\[\033[0;31m\]\w\[\033[0;34m\]\$(__git_ps1)\[\e[0m\] \[\033[0;32m\](\$(venv_ps1 %s)) \[\e[0m\]\$ "`, opts.Environment)
	}
}
