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
	ID api.ID
}

// DomainService orchestrates the creation of a hosted zone
type DomainService interface {
	CreatePrimaryHostedZone(ctx context.Context, opts CreatePrimaryHostedZoneOpts) (*api.HostedZone, error)
}

// DomainAPI invokes the API
type DomainAPI interface {
	CreatePrimaryHostedZone(opts CreatePrimaryHostedZoneOpts) (*api.HostedZone, error)
}

// DomainStore stores the data
type DomainStore interface {
	SaveHostedZone(*HostedZone) (*store.Report, error)
	GetPrimaryHostedZone(id api.ID) (*HostedZone, error)
}
