package client

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// We are shadowing some interfaces for now, but
// this is probably not sustainable.

// ParameterService implements the business logic
type ParameterService interface {
	api.ParameterService
}

// ParameterAPI invokes REST API endpoints
type ParameterAPI interface {
	api.ParameterCloudProvider
}

// ParameterStore stores the state
type ParameterStore interface {
	SaveSecret(*api.SecretParameter) (*store.Report, error)
}
