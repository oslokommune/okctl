package api

import "context"

type Domain struct {
	Repository   string
	Environment  string
	FQDN         string
	HostedZoneID string
	NameServers  []string
}

type CreateDomainOpts struct {
	Repository  string
	Environment string
	FQDN        string
}

type DomainService interface {
	CreateDomain(ctx context.Context, opts CreateDomainOpts) (*Domain, error)
}

type DomainCloudProvider interface {
	CreateDomain(opts CreateDomainOpts) (*Domain, error)
}

type DomainStore interface {
	SaveDomain(*Domain) error
}
