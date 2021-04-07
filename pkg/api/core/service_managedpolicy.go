package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type managedPolicyService struct {
	provider api.ManagedPolicyCloudProvider
}

func (m *managedPolicyService) CreatePolicy(_ context.Context, opts api.CreatePolicyOpts) (*api.ManagedPolicy, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	p, err := m.provider.CreatePolicy(opts)
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

func (m *managedPolicyService) CreateBlockstoragePolicy(_ context.Context, opts api.CreateBlockstoragePolicy) (*api.ManagedPolicy, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	got, err := m.provider.CreateBlockstoragePolicy(opts)
	if err != nil {
		return nil, errors.E(err, "creating blockstorage policy", errors.Internal)
	}

	return got, nil
}

func (m *managedPolicyService) DeleteBlockstoragePolicy(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = m.provider.DeleteBlockstoragePolicy(id)
	if err != nil {
		return errors.E(err, "deleting blockstorage policy", errors.Internal)
	}

	return nil
}

func (m *managedPolicyService) CreateAutoscalerPolicy(_ context.Context, opts api.CreateAutoscalerPolicy) (*api.ManagedPolicy, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	got, err := m.provider.CreateAutoscalerPolicy(opts)
	if err != nil {
		return nil, errors.E(err, "creating autoscaler policy", errors.Internal)
	}

	return got, nil
}

func (m *managedPolicyService) DeleteAutoscalerPolicy(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = m.provider.DeleteAutoscalerPolicy(id)
	if err != nil {
		return errors.E(err, "deleting autoscaler policy", errors.Internal)
	}

	return nil
}

func (m *managedPolicyService) DeleteExternalDNSPolicy(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errors.E(err, "failed to validate id", errors.Invalid)
	}

	err = m.provider.DeleteExternalDNSPolicy(id)
	if err != nil {
		return errors.E(err, "failed to delete external dns policy", errors.Internal)
	}

	return nil
}

func (m *managedPolicyService) CreateExternalDNSPolicy(_ context.Context, opts api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate external dns opts")
	}

	got, err := m.provider.CreateExternalDNSPolicy(opts)
	if err != nil {
		return nil, errors.E(err, "failed to create external dns policy")
	}

	return got, nil
}

// NewManagedPolicyService returns an initialised managed policy service
func NewManagedPolicyService(provider api.ManagedPolicyCloudProvider) api.ManagedPolicyService {
	return &managedPolicyService{
		provider: provider,
	}
}
