package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type managedPolicyService struct {
	provider api.ManagedPolicyCloudProvider
}

func (m *managedPolicyService) CreatePolicy(ctx context.Context, opts api.CreatePolicyOpts) (*api.ManagedPolicy, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	p, err := m.provider.CreatePolicy(ctx, opts)
	if err != nil {
		return nil, errors.E(err, "creating managed policy", errors.Internal)
	}

	return p, nil
}

func (m *managedPolicyService) DeletePolicy(_ context.Context, opts api.DeletePolicyOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = m.provider.DeletePolicy(opts)
	if err != nil {
		return errors.E(err, "deleting managed policy", errors.Internal)
	}

	return nil
}

// NewManagedPolicyService returns an initialised managed policy service
func NewManagedPolicyService(provider api.ManagedPolicyCloudProvider) api.ManagedPolicyService {
	return &managedPolicyService{
		provider: provider,
	}
}
