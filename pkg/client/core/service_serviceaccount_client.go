package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type serviceAccountService struct {
	service api.ServiceAccountService
	state   client.ServiceAccountState
}

func (m *serviceAccountService) CreateServiceAccount(context context.Context, opts client.CreateServiceAccountOpts) (*client.ServiceAccount, error) {
	s, err := m.service.CreateServiceAccount(context, api.CreateServiceAccountOpts{
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

	err = m.state.SaveServiceAccount(sa)
	if err != nil {
		return nil, err
	}

	return sa, nil
}

func (m *serviceAccountService) DeleteServiceAccount(context context.Context, opts client.DeleteServiceAccountOpts) error {
	err := m.service.DeleteServiceAccount(context, api.DeleteServiceAccountOpts{
		ID:     opts.ID,
		Name:   opts.Name,
		Config: opts.Config,
	})
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
	service api.ServiceAccountService,
	state client.ServiceAccountState,
) client.ServiceAccountService {
	return &serviceAccountService{
		service: service,
		state:   state,
	}
}
