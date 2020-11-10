package shellgetter

import (
	"bufio"
	"bytes"
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

type shell struct {
	Command   string
	ShellType ShellType
}

// shellGetter provides an interface for detecting which shell the user uses
type shellGetter interface {
	// Get detects which shell is used
	Get() (shell, error)
}

type loginShellGetter struct {
	envVars         map[string]string
	etcStorer       storage.Storer
	currentUsername string
}

func NewShellGetter(osEnvVars map[string]string, etcStorer storage.Storer, currentUsername string) *loginShellGetter {
	// TODO: Les OKCTL_SHELL. Hvis satt til noe fornuftig, lag en OkctlShellGetter. Hvis ikke, lag en loginShellGetter.
	return &loginShellGetter{
		envVars:         osEnvVars,
		etcStorer:       etcStorer,
		currentUsername: currentUsername,
	}
}

func (sd *loginShellGetter) Get() (shell, error) {
	cmd, err := sd.getShellCmd()
	if err != nil {
		return shell{}, fmt.Errorf("could not detect shell: %w", err)
	}

	var shellType ShellType
	if sd.shellIsBash(cmd) {
		shellType = ShellTypeBash
	} else if sd.shellIsZsh(cmd) {
		shellType = ShellTypeZsh
	} else {
		shellType = ShellTypeUnknown
	}

	return shell{
		Command:   cmd,
		ShellType: shellType,
	}, nil
}

// GetShellCmd detects which shell to run, and returns the command to run it. It detects the shell by looking for the
// environment variable OKCTL_SHELL. If it doesn't exist, it gets the user's default login shell from /etc/passwd.
func (sd *loginShellGetter) getShellCmd() (string, error) {
	shell, ok := sd.envVars["OKCTL_SHELL"]
	if ok {
		return shell, nil
	}

	line, err := sd.findUserLoginShell()
	if err != nil {
		return "", fmt.Errorf("could not get current user's login shell: %w", err)
	}

	shellCmd := sd.parseShellCmd(line)

	return shellCmd, nil
}

func (sd *loginShellGetter) findUserLoginShell() (string, error) {
	fileContents, err := sd.etcStorer.ReadAll("/passwd")
	if err != nil {
		return "", fmt.Errorf("failed to open /etc/passwd when detecting user login shell: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(fileContents))

	loginShell := ""
	for scanner.Scan() {
		loginShell = scanner.Text()

		if strings.Contains(scanner.Text(), sd.currentUsername) {
			return loginShell, nil
		}
	}

	return "", fmt.Errorf("failed to find '%s' in /etc/passwd", sd.currentUsername)
}

// TODO: Move to zsh and bash specific stuff

func (sd *loginShellGetter) parseShellCmd(line string) string {
	split := strings.Split(line, ":")
	return split[len(split)-1]
}

// ShellIsBash returns true if provided command will run bash
func (sd *loginShellGetter) shellIsBash(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "bash")
}

// ShellIsZsh returns true if provided command will run zsh
func (sd *loginShellGetter) shellIsZsh(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "zsh")
}
