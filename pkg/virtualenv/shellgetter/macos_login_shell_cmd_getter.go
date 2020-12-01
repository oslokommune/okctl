package shellgetter

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

type macOsLoginShellCmdGetter struct {
	cmdGetter MacOsUserShellCmdGetter
}

func newMacOsLoginShellCmdGetter(cmdGetter MacOsUserShellCmdGetter) *macOsLoginShellCmdGetter {
	return &macOsLoginShellCmdGetter{
		cmdGetter: cmdGetter,
	}
}

// Get detects which shell to run by getting the user's default login shell from /etc/passwd.
func (g *macOsLoginShellCmdGetter) Get() (string, error) {
	line, err := g.cmdGetter.RunDsclUserShell()
	if err != nil {
		return "", fmt.Errorf("could not get current user's login shell: %w", err)
	}

	shellCmd := g.parseShellCmd(line)

	return shellCmd, nil
}

// parseShellCmd transforms "UserShell: /bin/bash" to "/bin/bash"
func (*macOsLoginShellCmdGetter) parseShellCmd(line string) string {
	split := strings.Split(line, ":")
	afterColon := split[len(split)-1]

	return strings.TrimSpace(afterColon)
}

// MacOsUserShellCmdGetter gets the user's login shell on macOS
type MacOsUserShellCmdGetter interface {
	RunDsclUserShell() (string, error)
}

// NewMacOsCmdGetter returns a MacOsUserShellCmdGetter
func NewMacOsCmdGetter(userHomeDir string) MacOsUserShellCmdGetter {
	return &DefaultMacOsShellGetter{
		UserHomeDir: userHomeDir,
	}
}

// DefaultMacOsShellGetter is the default implementation for getting the user's login shell on macOS
type DefaultMacOsShellGetter struct {
	UserHomeDir string
}

// RunDsclUserShell returns the user's login shell as returned by dscl.
// Example: "UserShell: /bin/bash"
func (m *DefaultMacOsShellGetter) RunDsclUserShell() (string, error) {
	userHomeDir := path.Join(m.UserHomeDir) + "/"
	//nolint: gosec
	cmd := exec.Command("dscl", ".", "-read", userHomeDir, "UserShell")

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not run 'dscl' to get login shell: %w", err)
	}

	return string(out), nil
}
