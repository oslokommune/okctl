package cfn

import (
	"fmt"
)

// Not sure about this approach, will need to read up.

const (
	// DefaultStackNamePrefix defines a prefix added to all cloud
	// formation stacks
	DefaultStackNamePrefix = "okctl"
	// DefaultStackNameVpcID defines an identifier for a vpc
	DefaultStackNameVpcID = "vpc"
	// DefaultStackNameExternalSecretsPolicyID defines an identifier for an external secrets policy
	DefaultStackNameExternalSecretsPolicyID = "externalsecretspolicy"
)

// StackNamer knows how to name cloud formation stacks
type StackNamer struct{}

// NewStackNamer returns an initialised stack namer
func NewStackNamer() *StackNamer {
	return &StackNamer{}
}

// Vpc returns the stack name of a vpc
func (n *StackNamer) Vpc(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameVpcID,
		repository,
		env,
	)
}

// ExternalSecretsPolicy returns the stack name of an external secrets policy
func (n *StackNamer) ExternalSecretsPolicy(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameExternalSecretsPolicyID,
		repository,
		env,
	)
}
