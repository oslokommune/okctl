package client

import "github.com/oslokommune/okctl/pkg/api"

// We are shadowing the api interfaces for now, but
// this is probably not sustainable.

// DomainService orchestrates the creation of a hosted zone
type DomainService interface {
	api.DomainService
}

// DomainAPI invokes the API
type DomainAPI interface {
	api.DomainCloudProvider
}

// DomainStore stores the data
type DomainStore interface {
	api.DomainStore
}
