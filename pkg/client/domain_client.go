package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// HostedZone contains the state we are interested in
type HostedZone struct {
	ID                     api.ID
	IsDelegated            bool
	Primary                bool
	Managed                bool
	FQDN                   string
	Domain                 string
	HostedZoneID           string
	NameServers            []string
	StackName              string
	CloudFormationTemplate []byte
}

// CreatePrimaryHostedZoneOpts is the required inputs
type CreatePrimaryHostedZoneOpts struct {
	ID            api.ID
	Domain        string
	FQDN          string
	NameServerTTL int64
}

// DeletePrimaryHostedZoneOpts is the require inputs
type DeletePrimaryHostedZoneOpts struct {
	ID           api.ID
	HostedZoneID string
}

// DomainService orchestrates the creation of a hosted zone
type DomainService interface {
	CreatePrimaryHostedZone(ctx context.Context, opts CreatePrimaryHostedZoneOpts) (*HostedZone, error)
	GetPrimaryHostedZone(ctx context.Context) (*HostedZone, error)
	DeletePrimaryHostedZone(ctx context.Context, opts DeletePrimaryHostedZoneOpts) error
	SetHostedZoneDelegation(ctx context.Context, domain string, delegated bool) error
}

// DomainAPI invokes the API
type DomainAPI interface {
	CreatePrimaryHostedZone(opts CreatePrimaryHostedZoneOpts) (*HostedZone, error)
	DeletePrimaryHostedZone(domain string, opts DeletePrimaryHostedZoneOpts) error
}

// DomainStore stores the data
type DomainStore interface {
	SaveHostedZone(*HostedZone) (*store.Report, error)
	RemoveHostedZone(domain string) (*store.Report, error)
}

// DomainState implements the in-memory state handling
type DomainState interface {
	SaveHostedZone(zone *HostedZone) error
	UpdateHostedZone(zone *HostedZone) error
	RemoveHostedZone(domain string) error
	GetHostedZone(domain string) (*HostedZone, error)
	GetPrimaryHostedZone() (*HostedZone, error)
	GetHostedZones() ([]*HostedZone, error)
}
