package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type containerRepositoryDirectClient struct {
	service api.ContainerRepositoryService
}

func (c *containerRepositoryDirectClient) CreateContainerRepository(opts api.CreateContainerRepositoryOpts) (*api.ContainerRepository, error) {
	return c.service.CreateContainerRepository(context.Background(), &opts)
}

func (c *containerRepositoryDirectClient) DeleteContainerRepository(opts api.DeleteContainerRepositoryOpts) error {
	return c.service.DeleteContainerRepository(context.Background(), &opts)
}

// NewContainerRepositoryAPI returns an initialised REST API invoker
func NewContainerRepositoryAPI(service api.ContainerRepositoryService) client.ContainerRepositoryAPI {
	return &containerRepositoryDirectClient{
		service: service,
	}
}
