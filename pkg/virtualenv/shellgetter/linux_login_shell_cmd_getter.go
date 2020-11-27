package shellgetter

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"
)

type linuxLoginShellCmdGetter struct {
	etcStorer       storage.Storer
	currentUsername string
}

func newLinuxLoginShellCmdGetter(s storage.Storer, u string) shellCmdGetter {
	return &linuxLoginShellCmdGetter{
		etcStorer:       s,
		currentUsername: u,
	}
}

// Get detects which shell to run by getting the user's default login shell from /etc/passwd.
func (g *linuxLoginShellCmdGetter) Get() (string, error) {
	line, err := g.findUserLoginShell()
	if err != nil {
		return "", fmt.Errorf("could not get current user's login shell: %w", err)
	}

	shellCmd := g.parseShellCmd(line)

	return shellCmd, nil
}

func (g *linuxLoginShellCmdGetter) findUserLoginShell() (string, error) {
	fileContents, err := g.etcStorer.ReadAll("/passwd")
	if err != nil {
		return "", fmt.Errorf("failed to open /etc/passwd when detecting user login shell: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(fileContents))

	loginShell := ""
	for scanner.Scan() {
		loginShell = scanner.Text()

		if strings.Contains(scanner.Text(), g.currentUsername) {
			return loginShell, nil
		}
	}

	return "", fmt.Errorf("failed to find '%s' in /etc/passwd", g.currentUsername)
}

func (*linuxLoginShellCmdGetter) parseShellCmd(line string) string {
	split := strings.Split(line, ":")
	return split[len(split)-1]
}
