// Package okctlupgrade knows how to run an okctl upgrade
package okctlupgrade

import (
	"io"

	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/sirupsen/logrus"
)

const (
	// BinaryNameFormat is a string for formatting okctl upgrade binary names
	BinaryNameFormat = "okctl-upgrade_%s"
)

// BinaryRunner stores state for running the cli
type BinaryRunner struct {
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
) *BinaryRunner {
	return &BinaryRunner{
		repoDir:    repoDir,
		progress:   progress,
		logger:     logger,
		binaryPath: binaryPath,
		cmdFn:      cmdFn,
	}
}

// Run runs the okctl upgrade binary
func (u *BinaryRunner) Run(force bool) ([]byte, error) {
	var err error

	var envs []string

	var args []string

	if u.doDebug {
		args = append(args, "--debug")
	}

	if force {
		args = append(args, "--force")
	}

	runner := run.New(nil, u.repoDir, u.binaryPath, envs, u.cmdFn)

	output, err := runner.Run(u.progress, args)
	if err != nil {
		return nil, err
	}

	return output, err
}

// SetDebug sets whether we should increase log output from eksctl,
// the default behavior is off
func (u *BinaryRunner) SetDebug(enable bool) {
	u.doDebug = enable
}
