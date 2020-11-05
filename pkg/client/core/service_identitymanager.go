package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cognito"

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

// DeleteIdentityPool and all users
func (s *identityManagerService) DeleteIdentityPool(ctx context.Context, provider v1alpha1.CloudProvider, opts api.ID) error {
	err := s.spinner.Start("deleting identity-pool")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	pool := s.state.GetIdentityPool()

	c := cognito.New(provider)

	err = c.DeleteAuthDomain(pool.UserPoolID, pool.AuthDomain)
	if err != nil {
		return err
	}

	err = c.DeleteUserPool(pool.UserPoolID)
	if err != nil {
		return err
	}

	return nil
}

func (s *identityManagerService) CreateIdentityPoolUser(ctx context.Context, opts api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error) {
	err := s.spinner.Start("identity-pool-user")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	u, err := s.api.CreateIdentityPoolUser(opts)
	if err != nil {
		return nil, fmt.Errorf("creating identity pool user: %w", err)
	}

	r1, err := s.store.SaveIdentityPoolUser(u)
	if err != nil {
		return nil, fmt.Errorf("storing identity pool user: %w", err)
	}

	r2, err := s.state.SaveIdentityPoolUser(u)
	if err != nil {
		return nil, fmt.Errorf("updating identity pool user state: %w", err)
	}

	err = s.report.ReportIdentityPoolUser(u, []*store.Report{r1, r2})
	if err != nil {
		return nil, fmt.Errorf("reporting on identity pool user: %w", err)
	}

	return u, nil
}

func (s *identityManagerService) CreateIdentityPoolClient(_ context.Context, opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
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

func (s *identityManagerService) CreateIdentityPool(_ context.Context, opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
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
