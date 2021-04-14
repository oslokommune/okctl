package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// ExternalDNS contains state about an external dns deployment
type ExternalDNS struct {
	Policy         *ManagedPolicy
	ServiceAccount *ServiceAccount
	Kube           *ExternalDNSKube
}

// ExternalDNSKube contains the kubernetes data
type ExternalDNSKube struct {
	ID           api.ID
	HostedZoneID string
	DomainFilter string
	Manifests    map[string][]byte
}

// CreateExternalDNSOpts contains required inputs
type CreateExternalDNSOpts struct {
	ID           api.ID
	HostedZoneID string
	Domain       string
}

// ExternalDNSService is a business logic implementation
type ExternalDNSService interface {
	CreateExternalDNS(ctx context.Context, opts CreateExternalDNSOpts) (*ExternalDNS, error)
	DeleteExternalDNS(ctx context.Context, id api.ID) error
}

// ExternalDNSAPI implements the API invocation
type ExternalDNSAPI interface {
	CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.ExternalDNSKube, error)
}

// ExternalDNSState implements the persistence layer
type ExternalDNSState interface {
	SaveExternalDNS(dns *ExternalDNS) error
	GetExternalDNS() (*ExternalDNS, error)
	RemoveExternalDNS() error
}
