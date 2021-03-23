package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetComponentContainerRepository matches the REST API route
	TargetComponentContainerRepository = "components/containerrepository"
)

type containerRepositoryAPI struct {
	client *HTTPClient
}

func (c *containerRepositoryAPI) CreateContainerRepository(opts api.CreateContainerRepositoryOpts) (*api.ContainerRepository, error) {
	into := &api.ContainerRepository{}
	return into, c.client.DoPost(TargetComponentContainerRepository, &opts, into)
}

func (c *containerRepositoryAPI) DeleteContainerRepository(opts api.DeleteContainerRepositoryOpts) error {
	return c.client.DoDelete(TargetComponentContainerRepository, &opts)
}

// NewContainerRepositoryAPI returns an initialised REST API invoker
func NewContainerRepositoryAPI(client *HTTPClient) client.ContainerRepositoryAPI {
	return &containerRepositoryAPI{
		client: client,
	}
}
