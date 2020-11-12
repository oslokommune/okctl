// Package shellgetter implements functionality for deciding which shell is the primary shell of the provided
// environment.
package shellgetter

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"
)

// ShellType enumerates shells we recognize
type ShellType string

const (
	// ShellTypeBash is a constant that identifies the Bash shell
	ShellTypeBash ShellType = "bash"

	// ShellTypeZsh is a constant that identifies the Zsh shell
	ShellTypeZsh ShellType = "zsh"

	// ShellTypeUnknown is a constant that is identifies the case when an unknown shell is used
	ShellTypeUnknown ShellType = "unknown"
)

// ShellGetter is a provider for getting a shell executable based on some environment
type ShellGetter struct {
	osEnvVars       map[string]string
	etcStorer       storage.Storer
	currentUsername string
}

// Shell contains data about a shell executable (like bash)
type Shell struct {
	Command   string
	ShellType ShellType
}

// New creates a new ShellGetter
func New(osEnvVars map[string]string, etcStorer storage.Storer, currentUsername string) *ShellGetter {
	return &ShellGetter{
		osEnvVars:       osEnvVars,
		etcStorer:       etcStorer,
		currentUsername: currentUsername,
	}
}

// Get returns a new Shell
func (g *ShellGetter) Get() (*Shell, error) {
	shellCmdGetter := g.createShellCmdGetter()

	shellCmd, err := shellCmdGetter.Get()
	if err != nil {
		return &Shell{}, fmt.Errorf("could not get shell command: %w", err)
	}

	var shellType ShellType

	switch {
	case g.shellIsBash(shellCmd):
		shellType = ShellTypeBash
	case g.shellIsZsh(shellCmd):
		shellType = ShellTypeZsh
	default:
		shellType = ShellTypeUnknown
	}

	return &Shell{
		Command:   shellCmd,
		ShellType: shellType,
	}, nil
}

func (g *ShellGetter) shellIsBash(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "bash")
}

func (g *ShellGetter) shellIsZsh(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "zsh")
}
