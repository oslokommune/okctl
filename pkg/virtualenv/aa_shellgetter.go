package virtualenv

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

type Shell struct {
	Cmd       string
	ShellType ShellType
}

// ShellGetter provides an interface for detecting which shell the user uses
type ShellGetter interface {
	// Get detects which shell is used
	Get() (Shell, error)
}

type LoginShellGetter struct {
	storer      storage.Storer
	osLookupEnv OsLookupEnv
	currentUser string
}

// OsLookupEnv has the same functionality as os.LookupEnv
type OsLookupEnv func(key string) (string, bool)

func newShellDecider(osLookupEnv OsLookupEnv, storer storage.Storer, currentUser string) *LoginShellGetter {
	return &LoginShellGetter{
		osLookupEnv: osLookupEnv,
		storer:      storer,
		currentUser: currentUser,
	}
}

func (sd *LoginShellGetter) Get() (Shell, error) {
	cmd, err := sd.getShellCmd()
	if err != nil {
	    return Shell{}, fmt.Errorf("could not detect shell: %w", err)
	}

	var shellType ShellType
	if ShellIsBash(cmd) {
		shellType = ShellTypeBash
	} else if ShellIsZsh(cmd) {
		shellType = ShellTypeZsh
	} else {
		shellType = ShellTypeUnknown
	}

	return Shell{
		Cmd:       cmd,
		ShellType: shellType,
	}, nil
}

// GetShellCmd detects which shell to run, and returns the command to run it. It detects the shell by looking for the
// environment variable OKCTL_SHELL. If it doesn't exist, it gets the user's default login shell from /etc/passwd.
func (sd *LoginShellGetter) getShellCmd() (string, error) {
	shell, ok := sd.osLookupEnv("OKCTL_SHELL")
	if ok {
		return shell, nil
	}

	line, err := sd.findUserLoginShell()
	if err != nil {
		return "", fmt.Errorf("couldn't get shell: %w", err)
	}

	shellCmd := getShellCmd(line)

	return shellCmd, nil
}

func (sd *LoginShellGetter) findUserLoginShell() (string, error) {
	fileContents, err := sd.storer.ReadAll("/passwd")
	if err != nil {
		return "", fmt.Errorf("failed to open /etc/passwd when detecting user login shell: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(fileContents))

	loginShell := ""
	for scanner.Scan() {
		loginShell = scanner.Text()

		if strings.Contains(scanner.Text(), sd.currentUser) {
			return loginShell, nil
		}
	}

	return "", fmt.Errorf("failed to find '%s' in /etc/passwd", sd.currentUser)
}

func getShellCmd(line string) string {
	split := strings.Split(line, ":")
	return split[len(split)-1]
}
