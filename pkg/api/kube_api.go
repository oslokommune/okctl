package api

import (
	"context"

	"github.com/oslokommune/okctl/pkg/kube/manifests/storageclass"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ExternalSecretsKube is the state of an external secrets deployment
type ExternalSecretsKube struct {
	ID        ID
	Manifests map[string][]byte
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

// CreateExternalSecretsOpts contains the required inputs
type CreateExternalSecretsOpts struct {
	ID        ID
	Manifests []Manifest
}

// Validate the inputs
func (o CreateExternalSecretsOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Manifests, validation.Required),
	)
}

// Manifest represents a single external secret
type Manifest struct {
	Name        string
	Namespace   string
	Annotations map[string]string
	Labels      map[string]string
	Data        []Data
}

// Validate manifest
func (m Manifest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Namespace, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Data, validation.Required),
	)
}

// Data represents the items in the manifest
type Data struct {
	Key  string
	Name string
}

// Validate data
func (d Data) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.Name, validation.Required),
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

// KubeService provides kube deployment service layer
type KubeService interface {
	CreateExternalDNSKubeDeployment(ctx context.Context, opts CreateExternalDNSKubeDeploymentOpts) (*ExternalDNSKube, error)
	CreateExternalSecrets(ctx context.Context, opts CreateExternalSecretsOpts) (*ExternalSecretsKube, error)
	DeleteNamespace(ctx context.Context, opts DeleteNamespaceOpts) error
	CreateStorageClass(ctx context.Context, opts CreateStorageClassOpts) (*StorageClassKube, error)
}

// KubeRun provides kube deployment run layer
type KubeRun interface {
	CreateExternalDNSKubeDeployment(opts CreateExternalDNSKubeDeploymentOpts) (*ExternalDNSKube, error)
	CreateExternalSecrets(opts CreateExternalSecretsOpts) (*ExternalSecretsKube, error)
	DeleteNamespace(opts DeleteNamespaceOpts) error
	CreateStorageClass(opts CreateStorageClassOpts) (*StorageClassKube, error)
}

// KubeStore provides kube store layer
type KubeStore interface {
	SaveExternalDNSKubeDeployment(kube *ExternalDNSKube) error
	GetExternalDNSKubeDeployment() (*ExternalDNSKube, error)
	SaveExternalSecrets(kube *ExternalSecretsKube) error
}
