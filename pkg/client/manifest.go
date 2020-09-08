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
}

// ManifestAPI invokes the API
type ManifestAPI interface {
	CreateExternalSecret(opts CreateExternalSecretOpts) (*api.Kube, error)
}

// ManifestStore defines the storage layer
type ManifestStore interface {
	SaveExternalSecret(externalSecret *ExternalSecret) (*store.Report, error)
}
