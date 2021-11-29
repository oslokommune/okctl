package api

import (
	"context"

	"github.com/oslokommune/okctl/pkg/kube/manifests/storageclass"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ExternalSecretsKube is the state of an external secrets deployment
type ExternalSecretsKube struct {
	ID        ID
	Name      string
	Namespace string
	Content   []byte
}

// StorageClassKube is the state of a storage class manifest
type StorageClassKube struct {
	ID       ID
	Name     string
	Manifest []byte
}

// ExternalDNSKube is the state of an external dns deployment
type ExternalDNSKube struct {
	ID           ID
	HostedZoneID string
	DomainFilter string
	Manifests    map[string][]byte
}

// CreateExternalDNSKubeDeploymentOpts contains input options
type CreateExternalDNSKubeDeploymentOpts struct {
	ID           ID
	HostedZoneID string
	DomainFilter string
}

// Validate the input
func (o CreateExternalDNSKubeDeploymentOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.HostedZoneID, validation.Required),
		validation.Field(&o.DomainFilter, validation.Required),
	)
}

// DeleteExternalSecretsOpts contains the required inputs
type DeleteExternalSecretsOpts struct {
	ID        ID
	Manifests map[string]string
}

// Validate the provided inputs
func (o DeleteExternalSecretsOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Manifests, validation.Required),
	)
}

// CreateExternalSecretsOpts contains the required inputs
type CreateExternalSecretsOpts struct {
	ID       ID
	Manifest Manifest
}

// Validate the inputs
func (o CreateExternalSecretsOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Manifest, validation.Required),
	)
}

const (
	// BackendTypeSecretsManager for AWS SecretsManager
	BackendTypeSecretsManager = "secretsManager"
	// BackendTypeParameterStore for AWS Parameter Store
	BackendTypeParameterStore = "systemManager"
)

// ExternalSecretSpecTemplate represents the template attribute of an ExternalSecret
type ExternalSecretSpecTemplate struct {
	StringData map[string]interface{}
}

// Manifest represents a single external secret
type Manifest struct {
	Name        string
	Namespace   string
	Backend     string
	Annotations map[string]string
	Labels      map[string]string
	Data        []Data
	Template    ExternalSecretSpecTemplate
}

// Validate manifest
func (m Manifest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Namespace, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Backend, validation.Required),
		validation.Field(&m.Data, validation.Required),
	)
}

// Data represents the items in the manifest
type Data struct {
	Key      string
	Name     string
	Property string
}

// Validate data
func (d Data) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Name, validation.Required),
	)
}

// Namespace contains the data for a namespace
type Namespace struct {
	ID        ID
	Namespace string
	Labels    map[string]string
	Manifest  []byte
}

// CreateNamespaceOpts contains the required inputs
type CreateNamespaceOpts struct {
	ID        ID
	Namespace string
	Labels    map[string]string
}

// Validate the inputs
func (o CreateNamespaceOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
	)
}

// DeleteNamespaceOpts provides the inputs
type DeleteNamespaceOpts struct {
	ID        ID
	Namespace string
}

// Validate the inputs
func (o DeleteNamespaceOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Namespace, validation.Required),
	)
}

// CreateStorageClassOpts provides the inputs
type CreateStorageClassOpts struct {
	ID          ID
	Name        string
	Parameters  *storageclass.EBSParameters
	Annotations map[string]string
}

// Validate the inputs options
func (o CreateStorageClassOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Parameters, validation.Required),
	)
}

// DeleteConfigMapOpts contains the required inputs
type DeleteConfigMapOpts struct {
	ID        ID
	Name      string
	Namespace string
}

// Validate the provided inputs
func (o DeleteConfigMapOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
	)
}

// CreateConfigMapOpts contains the required inputs
type CreateConfigMapOpts struct {
	ID        ID
	Name      string
	Namespace string
	Data      map[string]string
	Labels    map[string]string
}

// Validate the inputs
func (o CreateConfigMapOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
		validation.Field(&o.Data, validation.Required),
	)
}

// ConfigMap is the state of a kubernetes configmap
type ConfigMap struct {
	ID        ID
	Name      string
	Namespace string
	Manifest  []byte
}

// ScaleDeploymentOpts provides required inputs
type ScaleDeploymentOpts struct {
	ID        ID
	Name      string
	Namespace string
	Replicas  int32
}

// Validate the provided inputs
func (o ScaleDeploymentOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
	)
}

// KubeService provides kube deployment service layer
type KubeService interface {
	CreateExternalDNSKubeDeployment(ctx context.Context, opts CreateExternalDNSKubeDeploymentOpts) (*ExternalDNSKube, error)
	DeleteNamespace(ctx context.Context, opts DeleteNamespaceOpts) error
	CreateStorageClass(ctx context.Context, opts CreateStorageClassOpts) (*StorageClassKube, error)
	CreateExternalSecrets(ctx context.Context, opts CreateExternalSecretsOpts) (*ExternalSecretsKube, error)
	DeleteExternalSecrets(ctx context.Context, opts DeleteExternalSecretsOpts) error
	CreateConfigMap(ctx context.Context, opts CreateConfigMapOpts) (*ConfigMap, error)
	DeleteConfigMap(ctx context.Context, opts DeleteConfigMapOpts) error
	ScaleDeployment(ctx context.Context, opts ScaleDeploymentOpts) error
	CreateNamespace(ctx context.Context, opts CreateNamespaceOpts) (*Namespace, error)
	DisableEarlyDEMUX(ctx context.Context, clusterID ID) error
}

// KubeRun provides kube deployment run layer
type KubeRun interface {
	CreateExternalDNSKubeDeployment(opts CreateExternalDNSKubeDeploymentOpts) (*ExternalDNSKube, error)
	DeleteNamespace(opts DeleteNamespaceOpts) error
	CreateStorageClass(opts CreateStorageClassOpts) (*StorageClassKube, error)
	CreateExternalSecrets(opts CreateExternalSecretsOpts) (*ExternalSecretsKube, error)
	DeleteExternalSecrets(opts DeleteExternalSecretsOpts) error
	CreateConfigMap(opts CreateConfigMapOpts) (*ConfigMap, error)
	DeleteConfigMap(opts DeleteConfigMapOpts) error
	ScaleDeployment(opts ScaleDeploymentOpts) error
	CreateNamespace(opts CreateNamespaceOpts) (*Namespace, error)
	DisableEarlyDEMUX(ctx context.Context, clusterID ID) error
}
