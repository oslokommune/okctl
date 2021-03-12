// Package secret knows how to build secrets manager secrets
package secret

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/secretsmanager"
)

// Secret contains all state required for building
// the cloud formation resource and outputs
type Secret struct {
	StoredName string
	Opts       Opts
}

// NamedOutputs returns the outputs commonly used by other stacks or components
func (s *Secret) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(s.Name(), s.Ref()).NamedOutputs()
}

// Name returns the name of the cloud formation resource
func (s *Secret) Name() string {
	return s.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (s *Secret) Ref() string {
	return cloudformation.Ref(s.Name())
}

const (
	passwordLength = 16
)

// Resource returns the cloud formation resource for a SecretsManager secret
func (s *Secret) Resource() cloudformation.Resource {
	return &secretsmanager.Secret{
		GenerateSecretString: &secretsmanager.Secret_GenerateSecretString{
			ExcludeCharacters:    s.Opts.ExcludeCharacters,
			GenerateStringKey:    s.Opts.GenerateStringKey,
			PasswordLength:       passwordLength,
			SecretStringTemplate: s.Opts.SecretStringTemplate,
		},
		Name: s.Opts.FriendlyName,
	}
}

// Opts contains the required inputs
type Opts struct {
	FriendlyName         string
	ExcludeCharacters    string
	SecretStringTemplate string
	GenerateStringKey    string
}

// New returns an initialised cloud formation secret
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-secret.html
func New(resourceName string, opts Opts) *Secret {
	return &Secret{
		StoredName: resourceName,
		Opts:       opts,
	}
}

// NewRDSInstanceSecret returns an initialised cloud formation secret
// compatible with RDSInstance
func NewRDSInstanceSecret(resourceName, friendlyName, userName string) *Secret {
	return New(resourceName, Opts{
		FriendlyName:         friendlyName,
		ExcludeCharacters:    `"@/\`,
		SecretStringTemplate: fmt.Sprintf(`{"username": "%s"}`, userName),
		GenerateStringKey:    "password",
	})
}
