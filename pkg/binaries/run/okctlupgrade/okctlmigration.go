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
	Progress   io.Writer
	Logger     *logrus.Logger
	BinaryPath string
	WorkingDir string
	CmdFn      run.CmdFn
}

// New creates a new kubectl cli wrapper
func New(
	logger *logrus.Logger,
	progress io.Writer,
	binaryPath string,
	fn run.CmdFn,
) *OkctlUpgrade {
	return &OkctlUpgrade{
		Progress:   progress,
		Logger:     logger,
		BinaryPath: binaryPath,
		CmdFn:      fn,
	}
}

func (u OkctlUpgrade) Run() {
	// TODO implement like eksctl run.
	fmt.Printf("Dummy running %s. Done!", u.BinaryPath)
}
