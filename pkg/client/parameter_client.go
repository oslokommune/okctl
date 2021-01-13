package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

// We are shadowing some interfaces for now, but
// this is probably not sustainable.

// ParameterService implements the business logic
type ParameterService interface {
	api.ParameterService
	DeleteAllsecrets(ctx context.Context, cluster state.Cluster) error
}

// ParameterAPI invokes REST API endpoints
type ParameterAPI interface {
	api.ParameterCloudProvider
}

// ParameterStore stores the state
type ParameterStore interface {
	SaveSecret(*api.SecretParameter) (*store.Report, error)
}

// ParameterReport defines the reporting layer
type ParameterReport interface {
	SaveSecret(parameter *api.SecretParameter, report *store.Report) error
}
