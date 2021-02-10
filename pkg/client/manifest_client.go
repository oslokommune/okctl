package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// ExternalSecret is the content of a kubernetes external secret
type ExternalSecret struct {
	ID        api.ID
	Manifests map[string][]byte
}

// CreateExternalSecretOpts contains the required inputs
type CreateExternalSecretOpts struct {
	ID        api.ID
	Manifests []api.Manifest
}

// ManifestService implements the business logic
type ManifestService interface {
	CreateExternalSecret(ctx context.Context, opts CreateExternalSecretOpts) (*ExternalSecret, error)
	// I don't like this, but I think ManifestService is an equally big problem
	// so I will leave it here for now. We should refactor this functionality
	// somehow..
	DeleteNamespace(ctx context.Context, opts api.DeleteNamespaceOpts) error
}

// ManifestAPI invokes the API
type ManifestAPI interface {
	CreateExternalSecret(opts CreateExternalSecretOpts) (*api.ExternalSecretsKube, error)
	DeleteNamespace(opts api.DeleteNamespaceOpts) error
}

// ManifestStore defines the storage layer
type ManifestStore interface {
	SaveExternalSecret(externalSecret *ExternalSecret) (*store.Report, error)
}

// ManifestReport defines the report layer
type ManifestReport interface {
	SaveExternalSecret(secret *ExternalSecret, report *store.Report) error
}
