package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// DomainService orchestrates the creation of a hosted zone
type DomainService interface {
	CreateHostedZone(ctx context.Context, opts api.CreateHostedZoneOpts) (*api.HostedZone, error)
}

// DomainAPI invokes the API
type DomainAPI interface {
	CreateHostedZone(opts api.CreateHostedZoneOpts) (*api.HostedZone, error)
}

// DomainStore stores the data
type DomainStore interface {
	SaveHostedZone(*api.HostedZone) (*store.Report, error)
}
