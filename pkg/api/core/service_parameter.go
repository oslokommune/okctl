package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

type parameter struct {
	cloudProvider api.ParameterCloudProvider
	store         api.ParameterStore
}

// nolint: godox
// TODO, implement this in proper client / api fashion.
func (p *parameter) DeleteSecret(ctx context.Context, provider v1alpha1.CloudProvider, name string) error {
	panic("Unused on the api side. Pay off this tech debt.")
}

func (p *parameter) CreateSecret(ctx context.Context, opts api.CreateSecretOpts) (*api.SecretParameter, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate secret parameter input", errors.Invalid)
	}

	param, err := p.cloudProvider.CreateSecret(opts)
	if err != nil {
		return nil, errors.E(err, "failed to create secret parameter", errors.Internal)
	}

	err = p.store.SaveSecret(param)
	if err != nil {
		return nil, errors.E(err, "failed to store secret parameter", errors.Internal)
	}

	return param, nil
}

// NewParameterService returns an initialised parameter service
func NewParameterService(cloudProvider api.ParameterCloudProvider, store api.ParameterStore) api.ParameterService {
	return &parameter{
		cloudProvider: cloudProvider,
		store:         store,
	}
}
