package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Kube contains the state for a kube deployment
type Kube struct {
	HostedZoneID string
	DomainFilter string

	Manifests map[string][]byte
}

// CreateExternalDNSKubeDeploymentOpts contains input options
type CreateExternalDNSKubeDeploymentOpts struct {
	HostedZoneID string
	DomainFilter string
}

// Validate the input
func (o CreateExternalDNSKubeDeploymentOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.HostedZoneID, validation.Required),
		validation.Field(&o.DomainFilter, validation.Required),
	)
}

// CreateExternalSecretsOpts contains the required inputs
type CreateExternalSecretsOpts struct {
	Manifests []Manifest
}

// Validate the inputs
func (o CreateExternalSecretsOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Manifests, validation.Required),
	)
}

// Manifest represents a single external secret
type Manifest struct {
	Name      string
	Namespace string
	Data      []Data
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

// KubeService provides kube deployment service layer
type KubeService interface {
	CreateExternalDNSKubeDeployment(ctx context.Context, opts CreateExternalDNSKubeDeploymentOpts) (*Kube, error)
	CreateExternalSecrets(ctx context.Context, opts CreateExternalSecretsOpts) (*Kube, error)
}

// KubeRun provides kube deployment run layer
type KubeRun interface {
	CreateExternalDNSKubeDeployment(opts CreateExternalDNSKubeDeploymentOpts) (*Kube, error)
	CreateExternalSecrets(opts CreateExternalSecretsOpts) (*Kube, error)
}

// KubeStore provides kube store layer
type KubeStore interface {
	SaveExternalDNSKubeDeployment(kube *Kube) error
	GetExternalDNSKubeDeployment() (*Kube, error)
	SaveExternalSecrets(kube *Kube) error
}
