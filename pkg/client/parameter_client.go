package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// SecretParameter contains the state for a secret parameter
type SecretParameter struct {
	ID      api.ID
	Name    string
	Path    string
	Version int64
	Content string
}

// CreateSecretOpts contains the input required for creating a secret parameter
type CreateSecretOpts struct {
	ID     api.ID
	Name   string
	Secret string
}

// DeleteSecretOpts contains the input required for deleting a secret parameter
type DeleteSecretOpts struct {
	ID   api.ID
	Name string
}

// ParameterService implements the business logic
type ParameterService interface {
	CreateSecret(ctx context.Context, opts CreateSecretOpts) (*SecretParameter, error)
	DeleteSecret(ctx context.Context, opts DeleteSecretOpts) error
}

// ParameterAPI invokes REST API endpoints
type ParameterAPI interface {
	CreateSecret(opts api.CreateSecretOpts) (*api.SecretParameter, error)
	DeleteSecret(opts api.DeleteSecretOpts) error
}

// ParameterState stores the state
type ParameterState interface {
	SaveSecret(parameter *SecretParameter) error
	GetSecret(name string) (*SecretParameter, error)
	RemoveSecret(name string) error
}
