package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// ExternalSecrets is the content of an external-secrets deployment
type ExternalSecrets struct {
	Policy         *ManagedPolicy
	ServiceAccount *ServiceAccount
	Chart          *Helm
}

// CreateExternalSecretsOpts contains the required inputs
type CreateExternalSecretsOpts struct {
	ID api.ID
}

// ExternalSecretsService is an implementation of the business logic
type ExternalSecretsService interface {
	CreateExternalSecrets(ctx context.Context, opts CreateExternalSecretsOpts) (*ExternalSecrets, error)
	DeleteExternalSecrets(ctx context.Context, id api.ID) error
}

// ExternalSecretsState defines the state layer
type ExternalSecretsState interface {
	HasExternalSecrets() (bool, error)
}
