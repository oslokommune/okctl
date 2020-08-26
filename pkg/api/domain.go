package api

import "context"

// Domain contains the state after creating a domain
type Domain struct {
	Repository   string
	Environment  string
	FQDN         string
	Domain       string
	HostedZoneID string
	NameServers  []string
}

// CreateDomainOpts contains the input required for creating a domain
type CreateDomainOpts struct {
	Repository  string
	Environment string
	Domain      string
	FQDN        string
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
