package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// HostedZone contains the state we are interested in
type HostedZone struct {
	IsDelegated bool
	Primary     bool
	HostedZone  *api.HostedZone
}

// CreatePrimaryHostedZoneOpts is the required inputs
type CreatePrimaryHostedZoneOpts struct {
	ID     api.ID
	Domain string
	FQDN   string
}

// DomainService orchestrates the creation of a hosted zone
type DomainService interface {
	CreatePrimaryHostedZone(ctx context.Context, opts CreatePrimaryHostedZoneOpts) (*HostedZone, error)
}

// DomainAPI invokes the API
type DomainAPI interface {
	CreatePrimaryHostedZone(opts CreatePrimaryHostedZoneOpts) (*HostedZone, error)
}

// DomainStore stores the data
type DomainStore interface {
	SaveHostedZone(*HostedZone) (*store.Report, error)
	GetHostedZone(domain string) (*HostedZone, error)
}

// DomainState implements the in-memory state handling
type DomainState interface {
	SaveHostedZone(zone *HostedZone) (*store.Report, error)
	GetHostedZones() []*HostedZone
}

// DomainReport implements the report layer
type DomainReport interface {
	ReportCreatePrimaryHostedZone(zone *HostedZone, reports []*store.Report) error
}
