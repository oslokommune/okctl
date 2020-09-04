package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Domain contains the state after creating a domain
type Domain struct {
	ID                     ID
	FQDN                   string
	Domain                 string
	HostedZoneID           string
	NameServers            []string
	StackName              string
	CloudFormationTemplate []byte
}

// CreateDomainOpts contains the input required for creating a domain
type CreateDomainOpts struct {
	ID     ID
	Domain string
	FQDN   string
}

// Validate the inputs
func (o CreateDomainOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Domain, validation.Required),
		validation.Field(&o.FQDN, validation.Required),
	)
}

// DomainService provides the service layer
type DomainService interface {
	CreateDomain(ctx context.Context, opts CreateDomainOpts) (*Domain, error)
}

// DomainCloudProvider provides the cloud provider layer
type DomainCloudProvider interface {
	CreateDomain(opts CreateDomainOpts) (*Domain, error)
}

// DomainStore provides the storage layer
type DomainStore interface {
	SaveDomain(*Domain) error
}
