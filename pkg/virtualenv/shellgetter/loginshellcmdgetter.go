package shellgetter

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"
)

type loginShellCmdGetter struct {
	envVars         map[string]string
	etcStorer       storage.Storer
	currentUsername string
}

// GetShellCmd detects which shell to run by getting the user's default login shell from /etc/passwd.
func (g *loginShellCmdGetter) Get() (string, error) {
	line, err := g.findUserLoginShell()
	if err != nil {
		return "", fmt.Errorf("could not get current user's login shell: %w", err)
	}

	shellCmd := g.parseShellCmd(line)

	return shellCmd, nil
}

func (sd *loginShellCmdGetter) findUserLoginShell() (string, error) {
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

func (sd *loginShellCmdGetter) parseShellCmd(line string) string {
	split := strings.Split(line, ":")
	return split[len(split)-1]
}
