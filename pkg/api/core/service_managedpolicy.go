package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type managedPolicyService struct {
	provider api.ManagedPolicyCloudProvider
	store    api.ManagedPolicyStore
}

func (m *managedPolicyService) CreateExternalSecretsPolicy(_ context.Context, opts api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate create external secrets policy options", errors.Invalid)
	}

	got, err := m.provider.CreateExternalSecretsPolicy(opts)
	if err != nil {
		return nil, errors.E(err, "failed to create external secrets policy")
	}

	err = m.store.SaveExternalSecretsPolicy(got)
	if err != nil {
		return nil, errors.E(err, "failed to save external secrets policy")
	}

	return got, nil
}

// NewManagedPolicyService returns an initialised managed policy service
func NewManagedPolicyService(provider api.ManagedPolicyCloudProvider, store api.ManagedPolicyStore) api.ManagedPolicyService {
	return &managedPolicyService{
		provider: provider,
		store:    store,
	}
}
