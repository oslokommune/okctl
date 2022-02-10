// Package kubens provides a convenient way of running kubens commands
package kubens

const (
	// Name sets the name of the binary/cli
	Name = "kubens"
	// Version sets the currently used version of the binary/cli
	Version = "0.9.4"
)

// Kubens stores state for running the cli
type Kubens struct {
	BinaryPath string
}

// New creates a new kubens cli wrapper
func New(binaryPath string) *Kubens {
	return &Kubens{
		BinaryPath: binaryPath,
	}
}
