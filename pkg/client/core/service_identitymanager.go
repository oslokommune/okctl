package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type identityManagerService struct {
	spinner spinner.Spinner
	api     client.IdentityManagerAPI
	store   client.IdentityManagerStore
	state   client.IdentityManagerState
	report  client.IdentityManagerReport
}

func (s identityManagerService) CreateIdentityPoolClient(_ context.Context, opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	err := s.spinner.Start("identity-pool-client")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	c, err := s.api.CreateIdentityPoolClient(opts)
	if err != nil {
		return nil, fmt.Errorf("creating identity pool client: %w", err)
	}

	r1, err := s.store.SaveIdentityPoolClient(c)
	if err != nil {
		return nil, fmt.Errorf("storing identity pool client: %w", err)
	}

	r2, err := s.state.SaveIdentityPoolClient(c)
	if err != nil {
		return nil, fmt.Errorf("updating identity pool client state: %w", err)
	}

	err = s.report.ReportIdentityPoolClient(c, []*store.Report{r1, r2})
	if err != nil {
		return nil, fmt.Errorf("reporting on identity pool client: %w", err)
	}

	return c, nil
}

func (s identityManagerService) CreateIdentityPool(_ context.Context, opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	err := s.spinner.Start("identity-pool")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

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
	spinner spinner.Spinner,
	api client.IdentityManagerAPI,
	store client.IdentityManagerStore,
	state client.IdentityManagerState,
	report client.IdentityManagerReport,
) client.IdentityManagerService {
	return &identityManagerService{
		spinner: spinner,
		api:     api,
		store:   store,
		state:   state,
		report:  report,
	}
}
