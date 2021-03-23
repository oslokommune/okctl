package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type containerRepositoryService struct {
	provider api.ContainerRepositoryCloudProvider
}

func (c *containerRepositoryService) CreateContainerRepository(_ context.Context, opts *api.CreateContainerRepositoryOpts) (*api.ContainerRepository, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "invalid inputs", errors.Invalid)
	}

	repository, err := c.provider.CreateContainerRepository(opts)
	if err != nil {
		return nil, errors.E(err, "creating container repository", errors.Internal)
	}

	return repository, nil
}

func (c *containerRepositoryService) DeleteContainerRepository(_ context.Context, opts *api.DeleteContainerRepositoryOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "invalid inputs", errors.Invalid)
	}

	err = c.provider.DeleteContainerRepository(opts)
	if err != nil {
		return errors.E(err, "deleting container repository", errors.Internal)
	}

	return nil
}

// NewComponentService returns an initialised component service
func NewContainerRepositoryService(provider api.ContainerRepositoryCloudProvider) api.ContainerRepositoryService {
	return &containerRepositoryService{
		provider: provider,
	}
}
