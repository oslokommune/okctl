package binary

import (
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/clients/terraform"
	"github.com/oslokommune/okctl/pkg/credentials"
)

type client struct {
	binaryProvider      binaries.Provider
	credentialsProvider credentials.Provider
	version             string
}

// Initialize knows how to terraform init a directory
func (c client) Initialize(directory string) error {
	return c.runCommand("terraform/init", directory, []string{"init"})
}

// Apply knows how to terraform apply a directory
func (c client) Apply(directory string) error {
	return c.runCommand("terraform/apply", directory, []string{"apply", "-auto-approve"})
}

// New initializes a new Terraform binary client
func New(binaryProvider binaries.Provider, credentialsProvider credentials.Provider, version string) terraform.Client {
	return &client{
		binaryProvider:      binaryProvider,
		credentialsProvider: credentialsProvider,
		version:             version,
	}
}
