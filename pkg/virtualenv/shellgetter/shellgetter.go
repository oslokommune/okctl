package shellgetter

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/storage"
	"strings"
)

// ShellType enumerates shells we recognize
type ShellType string

const (
	// ShellTypeBash is a constant that identifies the Bash shell
	ShellTypeBash ShellType = "bash"

	// ShellTypeBash is a constant that identifies the Zsh shell
	ShellTypeZsh ShellType = "zsh"

	// ShellTypeBash is a constant that is identifies the case when an unknown shell is used
	ShellTypeUnknown ShellType = "unknown"
)

type ShellGetter struct {
	osEnvVars       map[string]string
	etcStorer       storage.Storer
	currentUsername string
}

type Shell struct {
	Command   string
	ShellType ShellType
}

type shellCmdGetter interface {
	// GetCmd returns a shell command based on the user's environment
	Get() (string, error)
}

func New(osEnvVars map[string]string, etcStorer storage.Storer, currentUsername string) *ShellGetter {
	return &ShellGetter{
		osEnvVars:       osEnvVars,
		etcStorer:       etcStorer,
		currentUsername: currentUsername,
	}
}

func (g *ShellGetter) Get() (*Shell, error) {
	shellCmdGetter := g.createShellCmdGetter()

	shellCmd, err := shellCmdGetter.Get()
	if err != nil {
		return &Shell{}, fmt.Errorf("could not get shell command: %w", err)
	}

	var shellType ShellType
	if g.shellIsBash(shellCmd) {
		shellType = ShellTypeBash
	} else if g.shellIsZsh(shellCmd) {
		shellType = ShellTypeZsh
	} else {
		shellType = ShellTypeUnknown
	}

	return &Shell{
		Command:   shellCmd,
		ShellType: shellType,
	}, nil
}

func (g *ShellGetter) createShellCmdGetter() shellCmdGetter {
	shellCmd, isSet := g.osEnvVars["OKCTL_SHELL"]

	if isSet {
		return &envShellCmdGetter{
			shellCmd: shellCmd,
		}
	} else {
		return &loginShellCmdGetter{
			envVars:         g.osEnvVars,
			etcStorer:       g.etcStorer,
			currentUsername: g.currentUsername,
		}
	}
}

func (g *ShellGetter) shellIsBash(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "bash")
}

func (g *ShellGetter) shellIsZsh(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "zsh")
}
