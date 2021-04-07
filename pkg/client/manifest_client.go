package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// ManifestType enumerates the available
// manifest types
type ManifestType string

// String returns the string representation
func (t ManifestType) String() string {
	return string(t)
}

// nolint: golint
const (
	ManifestTypeExternalSecret = "external-secret"
	ManifestTypeStorageClass   = "storage-class"
	ManifestTypeConfigMap      = "config-map"
	ManifestTypeNamespace      = "namespace"
)

// KubernetesManifest contains data about
// a manifest
type KubernetesManifest struct {
	ID        api.ID
	Name      string
	Namespace string
	Type      ManifestType
	Content   []byte
}

// CreateExternalSecretOpts contains the required inputs
type CreateExternalSecretOpts struct {
	ID        api.ID
	Name      string
	Namespace string
	Manifest  api.Manifest
}

// DeleteExternalSecretOpts contains the required inputs
type DeleteExternalSecretOpts struct {
	ID      api.ID
	Name    string
	Secrets map[string]string
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
	CreateStorageClass(ctx context.Context, opts api.CreateStorageClassOpts) (*KubernetesManifest, error)
	CreateExternalSecret(ctx context.Context, opts CreateExternalSecretOpts) (*KubernetesManifest, error)
	DeleteExternalSecret(ctx context.Context, opts DeleteExternalSecretOpts) error
	CreateConfigMap(ctx context.Context, opts CreateConfigMapOpts) (*KubernetesManifest, error)
	DeleteConfigMap(ctx context.Context, opts DeleteConfigMapOpts) error
	CreateNamespace(ctx context.Context, opts api.CreateNamespaceOpts) (*KubernetesManifest, error)
	DeleteNamespace(ctx context.Context, opts api.DeleteNamespaceOpts) error
	ScaleDeployment(ctx context.Context, opts api.ScaleDeploymentOpts) error
}

// ManifestAPI invokes the API
type ManifestAPI interface {
	CreateStorageClass(opts api.CreateStorageClassOpts) (*api.StorageClassKube, error)
	CreateExternalSecret(opts CreateExternalSecretOpts) (*api.ExternalSecretsKube, error)
	DeleteExternalSecret(opts api.DeleteExternalSecretsOpts) error
	CreateConfigMap(opts api.CreateConfigMapOpts) (*api.ConfigMap, error)
	DeleteConfigMap(opts api.DeleteConfigMapOpts) error
	CreateNamespace(opts api.CreateNamespaceOpts) (*api.Namespace, error)
	DeleteNamespace(opts api.DeleteNamespaceOpts) error
	ScaleDeployment(opts api.ScaleDeploymentOpts) error
}

// ManifestState defines the state layer
type ManifestState interface {
	SaveKubernetesManifests(manifests *KubernetesManifest) error
	GetKubernetesManifests(name string) (*KubernetesManifest, error)
	RemoveKubernetesManifests(name string) error
}
