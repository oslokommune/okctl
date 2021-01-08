package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type parameter struct {
	cloudProvider api.ParameterCloudProvider
	store         api.ParameterStore
}

func (p *parameter) DeleteSecret(ctx context.Context, opts api.DeleteSecretOpts) error {
	err := p.cloudProvider.DeleteSecret(api.DeleteSecretOpts{
		Name: opts.Name,
	})
	if err != nil {
		return err
	}

	return nil
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
