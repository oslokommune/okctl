package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// DomainService orchestrates the creation of a hosted zone
type DomainService interface {
	CreateDomain(ctx context.Context, opts api.CreateDomainOpts) (*api.Domain, error)
}

// DomainAPI invokes the API
type DomainAPI interface {
	CreateDomain(opts api.CreateDomainOpts) (*api.Domain, error)
}

// DomainStore stores the data
type DomainStore interface {
	SaveDomain(*api.Domain) (*store.Report, error)
}
