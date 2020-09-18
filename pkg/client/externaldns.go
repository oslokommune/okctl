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
	Kube           *api.ExternalDNSKube
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
	CreateExternalDNSPolicy(opts api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error)
	DeleteExternalDNSPolicy(id api.ID) error
	CreateExternalDNSServiceAccount(opts api.CreateExternalDNSServiceAccountOpts) (*api.ServiceAccount, error)
	DeleteExternalDNSServiceAccount(id api.ID) error
	CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.ExternalDNSKube, error)
}

// ExternalDNSStore implements the storage layer
type ExternalDNSStore interface {
	SaveExternalDNS(dns *ExternalDNS) (*store.Report, error)
	RemoveExternalDNS(id api.ID) (*store.Report, error)
}

// ExternalDNSReport implements the report layer
type ExternalDNSReport interface {
	ReportCreateExternalDNS(dns *ExternalDNS, report *store.Report) error
	ReportDeleteExternalDNS(report *store.Report) error
}
