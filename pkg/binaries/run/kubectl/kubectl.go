// Package kubectl provides a convenient way of running kubectl commands
package kubectl

const (
	// Name sets the name of the binary/cli
	Name = "kubectl"
	// Version sets the currently used version of the binary/cli
	Version = "1.16.8"
)

// Kubectl stores state for running the cli
type Kubectl struct {
	BinaryPath string
}

// New creates a new kubectl cli wrapper
func New(binaryPath string) *Kubectl {
	return &Kubectl{
		BinaryPath: binaryPath,
	}
}
