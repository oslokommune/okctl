package virtualenv

import (
	"bufio"
	"bytes"
	"fmt"
	"os/user"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"
)

// OsLookupEnv has the same functionality as os.LookupEnv
type OsLookupEnv func(key string) (string, bool)

// GetShellCmd detects which shell to run, and returns the command to run it. It detects the shell by looking for the
// environment variable OKCTL_SHELL. If it doesn't exist, it gets the user's default login shell from /etc/passwd.
func GetShellCmd(osLookupEnv OsLookupEnv, storer storage.Storer) (string, error) {
	shell, ok := osLookupEnv("OKCTL_SHELL")
	if ok {
		return shell, nil
	}

	line, err := findUserLoginShell(storer)
	if err != nil {
		return "", fmt.Errorf("couldn't get shell: %w", err)
	}

	shellCmd := getShellCmd(line)

	return shellCmd, nil
}

func findUserLoginShell(storer storage.Storer) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	fileContents, err := storer.ReadAll("/passwd")
	if err != nil {
		return "", fmt.Errorf("failed to open /etc/passwd when detecting user login shell: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(fileContents))

	loginShell := ""
	for scanner.Scan() {
		loginShell = scanner.Text()

		if strings.Contains(scanner.Text(), currentUser.Username) {
			return loginShell, nil
		}
	}

	return "", fmt.Errorf("failed to find '%s' in /etc/passwd", currentUser.Username)
}

func getShellCmd(line string) string {
	split := strings.Split(line, ":")
	return split[len(split)-1]
}
