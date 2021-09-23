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

// Flags contains the flags to pass to the binary when running it.
// As a safety measure for not running an upgade binary accidentally with --dry-run false, which is the default, a dry
// run flag is not included as a field here.
// Too see the semantics of these flags, see https://github.com/oslokommune/okctl-upgrade
type Flags struct {
	Debug   bool
	Confirm bool
}

// DryRun runs the okctl upgrade binary with the dry-run flag set to true
func (u *BinaryRunner) DryRun(flags Flags) ([]byte, error) {
	return u.doRun(flags, true)
}

// Run runs the okctl upgrade binary with the dry-run flag set to false
func (u *BinaryRunner) Run(flags Flags) ([]byte, error) {
	return u.doRun(flags, false)
}

func (u *BinaryRunner) doRun(flags Flags, dryRun bool) ([]byte, error) {
	var err error

	var envs []string

	var args []string

	if flags.Debug {
		args = append(args, "--debug")
	}

	if flags.Confirm {
		args = append(args, "--debug")
	}

	if dryRun {
		args = append(args, "--dry-run", "true")
	} else {
		args = append(args, "--dry-run", "false")
	}

	runner := run.New(nil, u.repoDir, u.binaryPath, envs, u.cmdFn)

	output, err := runner.Run(u.progress, args)
	if err != nil {
		return nil, err
	}

	return output, err
}
