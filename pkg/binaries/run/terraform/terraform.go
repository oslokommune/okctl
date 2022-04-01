// Package terraform provides information about the terraform binary
package terraform

const (
	// Name sets the name of the binary/cli
	Name = "terraform"
	// Version sets the currently used version of the binary/cli
	Version = "1.1.7"
)

// Terraform stores state for running the cli
type Terraform struct {
	BinaryPath string
}

// New creates a new Terraform cli wrapper
func New(binaryPath string) *Terraform {
	return &Terraform{
		BinaryPath: binaryPath,
	}
}
