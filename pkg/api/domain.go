package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// HostedZone contains the state for a hosted zone
type HostedZone struct {
	ID                     ID
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
}

// Validate the inputs
func (o CreateHostedZoneOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Domain, validation.Required),
		validation.Field(&o.FQDN, validation.Required),
	)
}

// DomainService provides the service layer
type DomainService interface {
	CreateHostedZone(ctx context.Context, opts CreateHostedZoneOpts) (*HostedZone, error)
}

// DomainCloudProvider provides the cloud provider layer
type DomainCloudProvider interface {
	CreateHostedZone(opts CreateHostedZoneOpts) (*HostedZone, error)
}

// DomainStore provides the storage layer
type DomainStore interface {
	SaveHostedZone(*HostedZone) error
}
