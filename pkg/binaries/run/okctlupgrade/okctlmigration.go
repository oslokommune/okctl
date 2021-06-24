// Package okctlupgrade knows how to run an okctl upgrade
package okctlupgrade

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/sirupsen/logrus"
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

// Run runs the okctl upgrade binary
func (u *OkctlUpgrade) Run() ([]byte, error) {
	var err error

	runner := u.runner()

	var args []string

	output, err := runner.Run(u.progress, args)
	if err != nil {
		return nil, err
	}

	return output, err
}

// Debug sets whether we should increase log output from eksctl,
// the default behavior is off
func (u *OkctlUpgrade) Debug(enable bool) {
	u.doDebug = enable
}

func (u *OkctlUpgrade) runner() run.Runner {
	var envs []string

	if u.doDebug {
		envs = append(envs, fmt.Sprintf("%s=true", context.DefaultDebugEnv))
	}

	// return run.New(u.logger, u.repoDir, u.binaryPath, envs, u.cmdFn), nil
	return run.New(nil, u.repoDir, u.binaryPath, envs, u.cmdFn)
}
