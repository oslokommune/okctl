package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Parameter contains the state for a parameter
type Parameter struct {
	AWSAccountID   string
	RepositoryName string
	Environment    string
	Name           string
	Path           string
	Version        int64
	Content        string
}

// SecretParameter contains the state for a secret parameter
type SecretParameter struct {
	Parameter
}

// CreateSecretOpts contains the input required for creating a secret parameter
type CreateSecretOpts struct {
	AWSAccountID   string
	RepositoryName string
	Environment    string
	Name           string
	Secret         string
}

// Validate the inputs
func (o CreateSecretOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.RepositoryName, validation.Required),
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Secret, validation.Required),
	)
}

// ParameterService defines the service layer operations
type ParameterService interface {
	CreateSecret(ctx context.Context, opts CreateSecretOpts) (*SecretParameter, error)
}

// ParameterCloudProvider defines the cloud layer operations
type ParameterCloudProvider interface {
	CreateSecret(opts CreateSecretOpts) (*SecretParameter, error)
}

// ParameterStore defines the storage operations
type ParameterStore interface {
	SaveSecret(*SecretParameter) error
}
