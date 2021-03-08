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

// ConfigMap is the content of a kubernetes configmap
type ConfigMap struct {
	ID        api.ID
	Name      string
	Namespace string
	Manifest  []byte
}

// Namespace is the content of a k8s namespace
type Namespace struct {
	ID        api.ID
	Namespace string
	Labels    map[string]string
	Manifest  []byte
}

// CreateConfigMapOpts contains the required inputs
type CreateConfigMapOpts struct {
	ID        api.ID
	Name      string
	Namespace string
	Data      map[string]string
	Labels    map[string]string
}

// DeleteConfigMapOpts contains the required inputs
type DeleteConfigMapOpts struct {
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
	CreateConfigMap(ctx context.Context, opts CreateConfigMapOpts) (*ConfigMap, error)
	DeleteConfigMap(ctx context.Context, opts DeleteConfigMapOpts) error
	ScaleDeployment(ctx context.Context, opts api.ScaleDeploymentOpts) error
	CreateNamespace(ctx context.Context, opts api.CreateNamespaceOpts) (*Namespace, error)
}

// ManifestAPI invokes the API
type ManifestAPI interface {
	DeleteNamespace(opts api.DeleteNamespaceOpts) error
	CreateStorageClass(opts api.CreateStorageClassOpts) (*api.StorageClassKube, error)
	CreateExternalSecret(opts CreateExternalSecretOpts) (*api.ExternalSecretsKube, error)
	DeleteExternalSecret(opts api.DeleteExternalSecretsOpts) error
	CreateConfigMap(opts api.CreateConfigMapOpts) (*api.ConfigMap, error)
	DeleteConfigMap(opts api.DeleteConfigMapOpts) error
	ScaleDeployment(opts api.ScaleDeploymentOpts) error
	CreateNamespace(opts api.CreateNamespaceOpts) (*api.Namespace, error)
}

// ManifestStore defines the storage layer
type ManifestStore interface {
	SaveStorageClass(sc *StorageClass) (*store.Report, error)
	SaveExternalSecret(externalSecret *ExternalSecret) (*store.Report, error)
	RemoveExternalSecret(secrets map[string]string) (*store.Report, error)
	SaveConfigMap(configMap *ConfigMap) (*store.Report, error)
	RemoveConfigMap(name, namespace string) (*store.Report, error)
	SaveNamespace(namespace *Namespace) (*store.Report, error)
	RemoveNamespace(namespace string) (*store.Report, error)
}

// ManifestReport defines the report layer
type ManifestReport interface {
	SaveStorageClass(sc *StorageClass, report *store.Report) error
	SaveExternalSecret(secret *ExternalSecret, report *store.Report) error
	RemoveExternalSecret(report *store.Report) error
	SaveConfigMap(secret *ConfigMap, report *store.Report) error
	RemoveConfigMap(report *store.Report) error
	SaveNamespace(namespace *Namespace, report *store.Report) error
	RemoveNamespace(namespace string, report *store.Report) error
}
