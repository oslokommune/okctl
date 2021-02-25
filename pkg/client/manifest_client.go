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

// DeleteExternalSecretOpts contains the required inputs
type DeleteExternalSecretOpts struct {
	ID      api.ID
	Secrets map[string]string
}

// StorageClass is the content of a kubernetes storage class
type StorageClass struct {
	ID       api.ID
	Name     string
	Manifest []byte
}

// ManifestService implements the business logic
// There is nothing inherently wrong with this service, but I think there
// exists opportunities to improve the way we interact with Kubernetes and apply
// and remove resources. That is basically what this service does, it handles
// kubernetes resources from the client side.
type ManifestService interface {
	CreateExternalSecret(ctx context.Context, opts CreateExternalSecretOpts) (*ExternalSecret, error)
	DeleteNamespace(ctx context.Context, opts api.DeleteNamespaceOpts) error
	CreateStorageClass(ctx context.Context, opts api.CreateStorageClassOpts) (*StorageClass, error)
	DeleteExternalSecret(ctx context.Context, opts DeleteExternalSecretOpts) error
}

// ManifestAPI invokes the API
type ManifestAPI interface {
	CreateExternalSecret(opts CreateExternalSecretOpts) (*api.ExternalSecretsKube, error)
	DeleteNamespace(opts api.DeleteNamespaceOpts) error
	CreateStorageClass(opts api.CreateStorageClassOpts) (*api.StorageClassKube, error)
	DeleteExternalSecret(opts api.DeleteExternalSecretsOpts) error
}

// ManifestStore defines the storage layer
type ManifestStore interface {
	SaveExternalSecret(externalSecret *ExternalSecret) (*store.Report, error)
	SaveStorageClass(sc *StorageClass) (*store.Report, error)
	RemoveExternalSecret(secrets map[string]string) (*store.Report, error)
}

// ManifestReport defines the report layer
type ManifestReport interface {
	SaveExternalSecret(secret *ExternalSecret, report *store.Report) error
	SaveStorageClass(sc *StorageClass, report *store.Report) error
	RemoveExternalSecret(report *store.Report) error
}
