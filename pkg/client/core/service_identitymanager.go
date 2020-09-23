package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type identityManagerService struct {
	api    client.IdentityManagerAPI
	store  client.IdentityManagerStore
	state  client.IdentityManagerState
	report client.IdentityManagerReport
}

func (s identityManagerService) CreateIdentityPool(_ context.Context, opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	pool, err := s.api.CreateIdentityPool(opts)
	if err != nil {
		return nil, fmt.Errorf("creating identity pool: %w", err)
	}

	r1, err := s.store.SaveIdentityPool(pool)
	if err != nil {
		return nil, fmt.Errorf("storing identity pool: %w", err)
	}

	r2, err := s.state.SaveIdentityPool(pool)
	if err != nil {
		return nil, fmt.Errorf("updating identity pool state: %w", err)
	}

	err = s.report.ReportIdentityPool(pool, []*store.Report{r1, r2})
	if err != nil {
		return nil, fmt.Errorf("reporting on identity pool: %w", err)
	}

	return pool, nil
}

// NewIdentityManagerService returns an initialised service
func NewIdentityManagerService(
	api client.IdentityManagerAPI,
	store client.IdentityManagerStore,
	state client.IdentityManagerState,
	report client.IdentityManagerReport,
) api.IdentityManagerService {
	return &identityManagerService{
		api:    api,
		store:  store,
		state:  state,
		report: report,
	}
}
