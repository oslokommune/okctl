package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Kube contains the state for a kube deployment
type Kube struct {
	HostedZoneID string
	DomainFilter string
	Manifests    map[string][]byte
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

// KubeService provides kube deployment service layer
type KubeService interface {
	CreateExternalDNSKubeDeployment(ctx context.Context, opts CreateExternalDNSKubeDeploymentOpts) (*Kube, error)
}

// KubeRun provides kube deployment run layer
type KubeRun interface {
	CreateExternalDNSKubeDeployment(opts CreateExternalDNSKubeDeploymentOpts) (*Kube, error)
}

// KubeStore provides kube store layer
type KubeStore interface {
	SaveExternalDNSKubeDeployment(kube *Kube) error
	GetExternalDNSKubeDeployment() (*Kube, error)
}
