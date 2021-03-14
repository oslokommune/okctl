package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type parameter struct {
	cloudProvider api.ParameterCloudProvider
}

func (p *parameter) DeleteSecret(_ context.Context, opts api.DeleteSecretOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = p.cloudProvider.DeleteSecret(opts)
	if err != nil {
		return errors.E(err, "deleting secret parameter", errors.Internal)
	}

	return nil
}

func (p *parameter) CreateSecret(_ context.Context, opts api.CreateSecretOpts) (*api.SecretParameter, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	param, err := p.cloudProvider.CreateSecret(opts)
	if err != nil {
		return nil, errors.E(err, "creating secret parameter", errors.Internal)
	}

	return param, nil
}

// NewParameterService returns an initialised parameter service
func NewParameterService(cloudProvider api.ParameterCloudProvider) api.ParameterService {
	return &parameter{
		cloudProvider: cloudProvider,
	}
}
