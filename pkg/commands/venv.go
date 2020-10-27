package commands

// GetShell detects which shell to run, and returns the command to run it.
func GetShell(osLookupEnv func(key string) (string, bool)) string {
	shell, ok := osLookupEnv("OKCTL_SHELL")
	if ok {
		return shell
	}

	shell, ok = osLookupEnv("SHELL")
	if ok {
		return shell
	}

	return "/bin/sh"
}
