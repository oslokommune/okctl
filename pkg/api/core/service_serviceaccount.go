// Package core implements the service layer
package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type serviceAccount struct {
	run api.ServiceAccountRun
}

func (c *serviceAccount) CreateServiceAccount(_ context.Context, opts api.CreateServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	err = c.run.CreateServiceAccount(opts.Config)
	if err != nil {
		return nil, errors.E(err, "creating service account", errors.Internal)
	}

	return &api.ServiceAccount{
		ID:        opts.ID,
		Name:      opts.Name,
		PolicyArn: opts.PolicyArn,
		Config:    opts.Config,
	}, nil
}

func (c *serviceAccount) DeleteServiceAccount(_ context.Context, opts api.DeleteServiceAccountOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = c.run.DeleteServiceAccount(opts.Config)
	if err != nil {
		return errors.E(err, "deleting service account", errors.Internal)
	}

	return nil
}

// NewServiceAccountService returns a service operator for the service account operations
func NewServiceAccountService(run api.ServiceAccountRun) api.ServiceAccountService {
	return &serviceAccount{
		run: run,
	}
}
