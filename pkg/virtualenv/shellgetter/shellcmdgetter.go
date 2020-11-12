package shellgetter

type shellCmdGetter interface {
	// GetCmd returns a shell command based on the user's environment
	Get() (string, error)
}

func (g *ShellGetter) createShellCmdGetter() shellCmdGetter {
	shellCmd, isSet := g.osEnvVars["OKCTL_SHELL"]

	if isSet {
		return &envShellCmdGetter{
			shellCmd: shellCmd,
		}
	}

	return &loginShellCmdGetter{
		envVars:         g.osEnvVars,
		etcStorer:       g.etcStorer,
		currentUsername: g.currentUsername,
	}
}
