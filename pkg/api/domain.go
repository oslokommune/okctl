package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Domain contains the state after creating a domain
type Domain struct {
	Repository             string
	Environment            string
	FQDN                   string
	Domain                 string
	HostedZoneID           string
	NameServers            []string
	StackName              string
	CloudFormationTemplate []byte
}

// CreateDomainOpts contains the input required for creating a domain
type CreateDomainOpts struct {
	Repository  string
	Environment string
	Domain      string
	FQDN        string
}

// Validate the inputs
func (o CreateDomainOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Domain, validation.Required),
		validation.Field(&o.Repository, validation.Required),
		validation.Field(&o.FQDN, validation.Required),
		validation.Field(&o.Environment, validation.Required),
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
