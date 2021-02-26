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

// NativeSecret is the content of a kubernetes secret
type NativeSecret struct {
	ID        api.ID
	Name      string
	Namespace string
	Manifest  []byte
}

// CreateNativeSecretOpts contains the required inputs
type CreateNativeSecretOpts struct {
	ID        api.ID
	Name      string
	Namespace string
	Data      map[string]string
	Labels    map[string]string
}

// DeleteNativeSecretOpts contains the required inputs
type DeleteNativeSecretOpts struct {
	ID        api.ID
	Name      string
	Namespace string
}

// ManifestService implements the business logic
// There is nothing inherently wrong with this service, but I think there
// exists opportunities to improve the way we interact with Kubernetes and apply
// and remove resources. That is basically what this service does, it handles
// kubernetes resources from the client side.
type ManifestService interface {
	DeleteNamespace(ctx context.Context, opts api.DeleteNamespaceOpts) error
	CreateStorageClass(ctx context.Context, opts api.CreateStorageClassOpts) (*StorageClass, error)
	CreateExternalSecret(ctx context.Context, opts CreateExternalSecretOpts) (*ExternalSecret, error)
	DeleteExternalSecret(ctx context.Context, opts DeleteExternalSecretOpts) error
	CreateNativeSecret(ctx context.Context, opts CreateNativeSecretOpts) (*NativeSecret, error)
	DeleteNativeSecret(ctx context.Context, opts DeleteNativeSecretOpts) error
}

// ManifestAPI invokes the API
type ManifestAPI interface {
	DeleteNamespace(opts api.DeleteNamespaceOpts) error
	CreateStorageClass(opts api.CreateStorageClassOpts) (*api.StorageClassKube, error)
	CreateExternalSecret(opts CreateExternalSecretOpts) (*api.ExternalSecretsKube, error)
	DeleteExternalSecret(opts api.DeleteExternalSecretsOpts) error
	CreateNativeSecret(opts api.CreateNativeSecretOpts) (*api.NativeSecret, error)
	DeleteNativeSecret(opts api.DeleteNativeSecretOpts) error
}

// ManifestStore defines the storage layer
type ManifestStore interface {
	SaveStorageClass(sc *StorageClass) (*store.Report, error)
	SaveExternalSecret(externalSecret *ExternalSecret) (*store.Report, error)
	RemoveExternalSecret(secrets map[string]string) (*store.Report, error)
	SaveNativeSecret(NativeSecret *NativeSecret) (*store.Report, error)
	RemoveNativeSecret(name, namespace string) (*store.Report, error)
}

// ManifestReport defines the report layer
type ManifestReport interface {
	SaveStorageClass(sc *StorageClass, report *store.Report) error
	SaveExternalSecret(secret *ExternalSecret, report *store.Report) error
	RemoveExternalSecret(report *store.Report) error
	SaveNativeSecret(secret *NativeSecret, report *store.Report) error
	RemoveNativeSecret(report *store.Report) error
}
