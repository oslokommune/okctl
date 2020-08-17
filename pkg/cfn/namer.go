package cfn

import (
	"fmt"
	"strings"
)

// Not sure about this approach, will need to read up.

const (
	// DefaultStackNamePrefix defines a prefix added to all cloud
	// formation stacks
	DefaultStackNamePrefix = "okctl"
	// DefaultStackNameVpcID defines an identifier for a vpc
	DefaultStackNameVpcID = "vpc"
)

// StackNamer knows how to name cloud formation stacks
type StackNamer struct{}

// NewStackNamer returns an initialised stack namer
func NewStackNamer() *StackNamer {
	return &StackNamer{}
}

// Vpc returns the stack name of a vpc
func (n *StackNamer) Vpc(repoName, env string) string {
	return fmt.Sprintf("%s-%s-%s%s",
		DefaultStackNamePrefix,
		DefaultStackNameVpcID,
		strings.Title(repoName),
		strings.Title(env),
	)
}
