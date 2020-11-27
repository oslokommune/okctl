// Package shellgetter implements functionality for deciding which shell is the primary shell of the provided
// environment.
package shellgetter

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/virtualenv/shelltype"
)

// ShellGetter is a provider for getting a shell executable based on some environment
type ShellGetter struct {
	Os                      Os
	macOsUserShellCmdGetter MacOsUserShellCmdGetter
	OsEnvVars               map[string]string
	EtcStorage              storage.Storer
	CurrentUsername         string
}

// NewShellGetter returns a new ShellGetter
func NewShellGetter(
	os Os,
	macOsUserShellGetter MacOsUserShellCmdGetter,
	osEnvVars map[string]string,
	etcStorage storage.Storer,
	currentUsername string) *ShellGetter {
	return &ShellGetter{
		Os:                      os,
		macOsUserShellCmdGetter: macOsUserShellGetter,
		OsEnvVars:               osEnvVars,
		EtcStorage:              etcStorage,
		CurrentUsername:         currentUsername,
	}
}

// Shell contains data about a shell executable (like bash)
type Shell struct {
	Command   string
	ShellType shelltype.ShellType
}

// Get returns a new Shell
func (g *ShellGetter) Get() (*Shell, error) {
	shellCmdGetter := g.createShellCmdGetter()

	shellCmd, err := shellCmdGetter.Get()
	if err != nil {
		return &Shell{}, fmt.Errorf("could not get shell command: %w", err)
	}

	var shellType shelltype.ShellType

	switch {
	case g.shellIsBash(shellCmd):
		shellType = shelltype.Bash
	case g.shellIsZsh(shellCmd):
		shellType = shelltype.Zsh
	default:
		shellType = shelltype.Unknown
	}

	return &Shell{
		Command:   shellCmd,
		ShellType: shellType,
	}, nil
}

func (g *ShellGetter) createShellCmdGetter() shellCmdGetter {
	shellCmd, isSet := g.OsEnvVars["OKCTL_SHELL"]

	if isSet {
		return &envShellCmdGetter{
			shellCmd: shellCmd,
		}
	}

	if g.Os == OsDarwin {
		return newMacOsLoginShellCmdGetter(g.macOsUserShellCmdGetter)
	}

	return newLinuxLoginShellCmdGetter(g.EtcStorage, g.CurrentUsername)
}

func (g *ShellGetter) shellIsBash(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "bash")
}

func (g *ShellGetter) shellIsZsh(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "zsh")
}
