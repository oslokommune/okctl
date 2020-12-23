package client

import (
	"context"

	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/config/state"

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

// DeletePrimaryHostedZoneOpts is the require inputs
type DeletePrimaryHostedZoneOpts struct {
	ID           api.ID
	HostedZoneID string
}

// DomainService orchestrates the creation of a hosted zone
type DomainService interface {
	CreatePrimaryHostedZone(ctx context.Context, opts CreatePrimaryHostedZoneOpts) (*HostedZone, error)
	CreatePrimaryHostedZoneWithoutUserinput(ctx context.Context, opts CreatePrimaryHostedZoneOpts) (*HostedZone, error)
	GetPrimaryHostedZone(ctx context.Context, id api.ID) (*HostedZone, error)
	DeletePrimaryHostedZone(ctx context.Context, opts DeletePrimaryHostedZoneOpts) error
}

// DomainAPI invokes the API
type DomainAPI interface {
	CreatePrimaryHostedZone(opts CreatePrimaryHostedZoneOpts) (*HostedZone, error)
	DeletePrimaryHostedZone(domain string, opts DeletePrimaryHostedZoneOpts) error
	DeleteHostedZoneRecords(provider v1alpha1.CloudProvider, hostedZoneID string) (*route53.ChangeResourceRecordSetsOutput, error)
}

// DomainStore stores the data
type DomainStore interface {
	SaveHostedZone(*HostedZone) (*store.Report, error)
	GetHostedZone(domain string) (*HostedZone, error)
	RemoveHostedZone(domain string) (*store.Report, error)
}

// DomainState implements the in-memory state handling
type DomainState interface {
	SaveHostedZone(zone *HostedZone) (*store.Report, error)
	GetHostedZones() []state.HostedZone
	RemoveHostedZone(domain string) (*store.Report, error)
}

// DomainReport implements the report layer
type DomainReport interface {
	ReportCreatePrimaryHostedZone(zone *HostedZone, reports []*store.Report) error
	ReportDeletePrimaryHostedZone(reports []*store.Report) error
}
