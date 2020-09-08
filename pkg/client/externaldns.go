package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// ExternalDNS contains state about an external dns deployment
type ExternalDNS struct {
	Policy         *api.ManagedPolicy
	ServiceAccount *api.ServiceAccount
	Kube           *api.Kube
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
}

// ExternalDNSAPI implements the API invocation
type ExternalDNSAPI interface {
	CreateExternalDNSPolicy(opts api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error)
	CreateExternalDNSServiceAccount(opts api.CreateExternalDNSServiceAccountOpts) (*api.ServiceAccount, error)
	CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.Kube, error)
}

// ExternalDNSStore implements the storage layer
type ExternalDNSStore interface {
	SaveExternalDNS(dns *ExternalDNS) (*store.Report, error)
}
