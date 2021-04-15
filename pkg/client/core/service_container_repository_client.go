package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type containerRepositoryService struct {
	api   client.ContainerRepositoryAPI
	state client.ContainerRepositoryState

	provider v1alpha1.CloudProvider
}

// CreateContainerRepository handles api, state, store and report orchistration for creation of a container repository
func (c *containerRepositoryService) CreateContainerRepository(_ context.Context, opts client.CreateContainerRepositoryOpts) (*client.ContainerRepository, error) {
	repository, err := c.api.CreateContainerRepository(api.CreateContainerRepositoryOpts{
		ClusterID: opts.ClusterID,
		Name:      opts.ImageName,
		StackName: cfn.NewStackNamer().ContainerRepository(opts.ImageName, opts.ClusterID.ClusterName),
	})
	if err != nil {
		return nil, err
	}

	containerRepository := &client.ContainerRepository{
		ClusterID:              repository.ClusterID,
		ImageName:              repository.Name,
		StackName:              repository.StackName,
		CloudFormationTemplate: repository.CloudFormationTemplate,
	}

	err = c.state.SaveContainerRepository(containerRepository)
	if err != nil {
		return nil, err
	}

	return containerRepository, nil
}

// DeleteContainerRepository handles api, state, store and report orchistration for deletion of a container repository
func (c *containerRepositoryService) DeleteContainerRepository(_ context.Context, opts client.DeleteContainerRepositoryOpts) error {
	err := c.api.DeleteContainerRepository(api.DeleteContainerRepositoryOpts{
		ClusterID: opts.ClusterID,
		StackName: cfn.NewStackNamer().ContainerRepository(opts.ImageName, opts.ClusterID.ClusterName),
	})
	if err != nil {
		return err
	}

	err = c.state.RemoveContainerRepository(opts.ImageName)
	if err != nil {
		return err
	}

	return nil
}

// NewContainerRepositoryService returns an initialised container repository service
func NewContainerRepositoryService(
	api client.ContainerRepositoryAPI,
	state client.ContainerRepositoryState,
	provider v1alpha1.CloudProvider,
) client.ContainerRepositoryService {
	return &containerRepositoryService{
		api:      api,
		state:    state,
		provider: provider,
	}
}
