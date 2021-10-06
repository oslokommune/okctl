package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// HostedZone contains the state for a hosted zone
type HostedZone struct {
	ID                     ID
	Managed                bool
	FQDN                   string
	Domain                 string
	HostedZoneID           string
	NameServers            []string
	StackName              string
	CloudFormationTemplate []byte
}

// CreateHostedZoneOpts contains required inputs
type CreateHostedZoneOpts struct {
	ID     ID
	Domain string
	FQDN   string
	NSTTL  int64
}

// Validate the inputs
func (o CreateHostedZoneOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Domain, validation.Required),
		validation.Field(&o.FQDN, validation.Required),
	)
}

// DeleteHostedZoneOpts contains required inputs
type DeleteHostedZoneOpts struct {
	ID           ID
	HostedZoneID string
	Domain       string
}

// Validate the inputs
func (o DeleteHostedZoneOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Domain, validation.Required),
	)
}

// DomainService provides the service layer
type DomainService interface {
	CreateHostedZone(ctx context.Context, opts CreateHostedZoneOpts) (*HostedZone, error)
	DeleteHostedZone(ctx context.Context, opts DeleteHostedZoneOpts) error
}

// DomainCloudProvider provides the cloud provider layer
type DomainCloudProvider interface {
	CreateHostedZone(ctx context.Context, opts CreateHostedZoneOpts) (*HostedZone, error)
	DeleteHostedZone(opts DeleteHostedZoneOpts) error
}
