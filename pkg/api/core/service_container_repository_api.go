package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type containerRepositoryService struct {
	provider api.ContainerRepositoryCloudProvider
}

// CreateContainerRepository handles creating a container repository
func (c *containerRepositoryService) CreateContainerRepository(_ context.Context, opts *api.CreateContainerRepositoryOpts) (*api.ContainerRepository, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating create opts", errors.Invalid)
	}

	repository, err := c.provider.CreateContainerRepository(opts)
	if err != nil {
		return nil, errors.E(err, "creating container repository", errors.Internal)
	}

	return repository, nil
}

// DeleteContainerRepository handles deleting a container repository
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

// EmptyContainerRepository removes all containers in a container repository
func (c *containerRepositoryService) EmptyContainerRepository(ctx context.Context, opts api.EmptyContainerRepositoryOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "invalid inputs", errors.Invalid)
	}

	err = c.provider.EmptyContainerRepository(ctx, opts)
	if err != nil {
		return errors.E(err, "emptying container repository", errors.Internal)
	}

	return nil
}

// NewContainerRepositoryService returns an initialised container repository service
func NewContainerRepositoryService(provider api.ContainerRepositoryCloudProvider) api.ContainerRepositoryService {
	return &containerRepositoryService{
		provider: provider,
	}
}
