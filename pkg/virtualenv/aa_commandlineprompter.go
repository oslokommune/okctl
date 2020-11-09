package virtualenv

// CommandLinePrompter provides an interface for configuring the command line prompt
type CommandLinePrompter interface {
	SetPrompt(ps1Dir string, venv *VirtualEnvironment) error
	Cleanup() error
}

func NewCommandLinePrompter(shellType ShellType) CommandLinePrompter {
	switch shellType {
	case ShellTypeBash:
		return &BashPrompter{}
	case ShellTypeZsh:
		return &ZshPrompter{}
	}

}




