package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type serviceAccountService struct {
	api   client.ServiceAccountAPI
	store client.ServiceAccountStore
	state client.ServiceAccountState
}

func (m *serviceAccountService) CreateServiceAccount(_ context.Context, opts client.CreateServiceAccountOpts) (*client.ServiceAccount, error) {
	s, err := m.api.CreateServiceAccount(api.CreateServiceAccountOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		PolicyArn: opts.PolicyArn,
		Config:    opts.Config,
	})
	if err != nil {
		return nil, err
	}

	sa := &client.ServiceAccount{
		ID:        s.ID,
		Name:      s.Name,
		PolicyArn: s.PolicyArn,
		Config:    s.Config,
	}

	err = m.store.SaveServiceAccount(sa)
	if err != nil {
		return nil, err
	}

	err = m.state.SaveServiceAccount(sa)
	if err != nil {
		return nil, err
	}

	return sa, nil
}

func (m *serviceAccountService) DeleteServiceAccount(_ context.Context, opts client.DeleteServiceAccountOpts) error {
	err := m.api.DeleteServiceAccount(api.DeleteServiceAccountOpts{
		ID:     opts.ID,
		Name:   opts.Name,
		Config: opts.Config,
	})
	if err != nil {
		return err
	}

	err = m.store.RemoveServiceAccount(opts.Name)
	if err != nil {
		return err
	}

	err = m.state.RemoveServiceAccount(opts.Name)
	if err != nil {
		return err
	}

	return nil
}

// NewServiceAccountService returns an initialised service
func NewServiceAccountService(
	api client.ServiceAccountAPI,
	store client.ServiceAccountStore,
	state client.ServiceAccountState,
) client.ServiceAccountService {
	return &serviceAccountService{
		api:   api,
		store: store,
		state: state,
	}
}
