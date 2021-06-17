package okctlupgrade

import "fmt"

const (
	// BinaryNameFormat is a string for formatting okctl upgrade binary names
	BinaryNameFormat = "okctl-upgrade_%s"
)

// OkctlUpgrade stores state for running the cli
type OkctlUpgrade struct {
	BinaryPath string
}

func (u OkctlUpgrade) Run() {
	// TODO implement like eksctl run
	fmt.Printf("Dummy running %s. Done!", u.BinaryPath)
}

// New creates a new kubectl cli wrapper
func New(binaryPath string) *OkctlUpgrade {
	return &OkctlUpgrade{
		BinaryPath: binaryPath,
	}
}
