package okctlupgrade

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/sirupsen/logrus"
	"io"
)

const (
	// BinaryNameFormat is a string for formatting okctl upgrade binary names
	BinaryNameFormat = "okctl-upgrade_%s"
)

// OkctlUpgrade stores state for running the cli
type OkctlUpgrade struct {
	doDebug    bool
	repoDir    string
	progress   io.Writer
	logger     *logrus.Logger
	binaryPath string
	cmdFn      run.CmdFn
}

// New creates a new okctl upgrade cli wrapper
func New(
	repoDir string,
	progress io.Writer,
	logger *logrus.Logger,
	binaryPath string,
	cmdFn run.CmdFn,
) *OkctlUpgrade {
	return &OkctlUpgrade{
		repoDir:    repoDir,
		progress:   progress,
		logger:     logger,
		binaryPath: binaryPath,
		cmdFn:      cmdFn,
	}
}

func (u *OkctlUpgrade) Run() ([]byte, error) {
	var err error

	runner, err := u.runner()
	if err != nil {
		return nil, fmt.Errorf("getting runner: %w", err)
	}

	var args []string

	return runner.Run(u.progress, args)
}

// Debug sets whether we should increase log output from eksctl,
// the default behavior is off
func (u *OkctlUpgrade) Debug(enable bool) {
	u.doDebug = enable
}

func (u *OkctlUpgrade) runner() (run.Runner, error) {
	var envs []string

	if u.doDebug {
		envs = append(envs, "OKCTL_DEBUG=true")
	}

	return run.New(u.logger, u.repoDir, u.binaryPath, envs, u.cmdFn), nil
}
