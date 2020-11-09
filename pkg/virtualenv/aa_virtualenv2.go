package virtualenv

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/storage"
)

type VirtualEnvCreator interface {
	Create(opts VirtualEnvironmentOpts) (*VirtualCommand, error)
}

type VirtualEnvCleaner interface {
	Clean() error
}

type VirtualEnver interface {
	VirtualEnvCreator
	VirtualEnvCleaner
}

// VirtualEnv contains state for maintaining a
type VirtualEnv struct {
	shellGetter ShellGetter // Maybe this shouldn't be stored in here, but rather just used to create the proper prompter
	prompt      CommandLinePrompter
}

type VirtualCommand struct {
	ShellCommand string
	Env 		 []string
}

// Ensure that VirtualEnv implementer the VirtualEnver interface
var _ VirtualEnver = &VirtualEnv{}

// Create the virtual env by orchestrating all the different parts
func (v *VirtualEnv) Create(opts VirtualEnvironmentOpts) (*VirtualCommand, error) {
	// 11111111111111 GET WHICH SHELL
	shell, err := v.shellGetter.Get()
	if err != nil {
		return &VirtualCommand{}, err
	}

	// 222222222222 GET THE ENVIRONMENT VARIABLES
	virtualCommand := &VirtualCommand{}

	venv, err := GetVirtualEnvironment(opts) // TODO: Don't deal with ps1 here
	if err != nil {
	    return nil, fmt.Errorf("couldn't get virtual environment: %w", err)
	}

	// Lag en command prompter, som tar shell.ShellType som input
	// command prompter tar som input:
	//   venv: modifiserer denne til å inneholde det som trengs på PATH (BASH_PROMPT eller ZDOTDIR)
	//   storage: modifiserer denne til å inneholdet som som trengs av eksekverbare filer (venv_ps1, og hvis zsh: .zshrc-fil)
	//




	ps1, overridePs1 := venv.env["OKCTL_PS1"]
	if overridePs1 {
		venv.env["PROMPT_COMMAND"] = fmt.Sprintf("PS1=%s", ps1)
	} else {
		venv.env["PROMPT_COMMAND"] = fmt.Sprintf(`PS1="\[\033[0;31m\]\w \[\033[0;34m\](\$(venv_ps1 %s)) \[\e[0m\]\$ "`, opts.Environment)
	}

	prompter := NewCommandLinePrompter(shell.ShellType)
	err = prompter.SetPrompt(opts.Ps1Dir, venv)
	if err != nil {
	    return nil, fmt.Errorf("could not build command line prompt: %w", err)
	}


	//
	//switch shell.ShellType {
	//case ShellTypeBash:
	//
	//	// create a bash prompter
	//	// 333333333333 UPDATE ENVIRONMENT VARIABLES WITH COMMANDT PROMPT STUFF
	//
	//case ShellTypeZsh:
	//	//virtualCommand.ShellCommand = "/bin/zsh"
	//
	//	// create a ZSH prompter
	//	// 333333333333 UPDATE ENVIRONMENT VARIABLES WITH COMMANDT PROMPT STUFF
	//default:
	//	return nil, fmt.Errorf("unknown shell type: %s", shell.ShellType)
	//}



	virtualCommand.ShellCommand = shell.Cmd
	virtualCommand.Env = toEnvVarsSlice(&venv.env)

	data, err := v.prompt.BuildPrompt(opts)
	if err != nil {
		return nil, fmt.Errorf("coudln't build prompt: %w", err)
	}

	err = v.prompt.SetPrompt(data, nil)
	if err != nil {
		return nil, fmt.Errorf("coudln't set prompt: %w", err)
	}

	return virtualCommand, nil
}

// Clean removes all remnants, if any such exist
func (v *VirtualEnv) Clean() error {
	return v.prompt.Cleanup()
}

// New returns a new virtual environment
func New(osLookupEnv OsLookupEnv, storer storage.Storer, currentUser string) (*VirtualEnv, error) {
	sg := newShellDecider(osLookupEnv, storer, currentUser)
	p := newBashPrompter()

	return &VirtualEnv{
		shellGetter: sg,
		prompt: p,
	}, nil
}
