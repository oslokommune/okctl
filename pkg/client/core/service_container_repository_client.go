package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type containerRepositoryService struct {
	service api.ContainerRepositoryService
	state   client.ContainerRepositoryState

	provider v1alpha1.CloudProvider
}

// CreateContainerRepository handles api, state, store and report orchistration for creation of a container repository
func (c *containerRepositoryService) CreateContainerRepository(context context.Context, opts client.CreateContainerRepositoryOpts) (*client.ContainerRepository, error) {
	repository, err := c.service.CreateContainerRepository(context, &api.CreateContainerRepositoryOpts{
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
		ApplicationName:        opts.ApplicationName,
		CloudFormationTemplate: repository.CloudFormationTemplate,
	}

	err = c.state.SaveContainerRepository(containerRepository)
	if err != nil {
		return nil, err
	}

	return containerRepository, nil
}

// DeleteContainerRepository handles api, state, store and report orchistration for deletion of a container repository
func (c *containerRepositoryService) DeleteContainerRepository(context context.Context, opts client.DeleteContainerRepositoryOpts) error {
	err := c.service.DeleteContainerRepository(context, &api.DeleteContainerRepositoryOpts{
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

// EmptyContainerRepository deletes all images present in a container repository
func (c *containerRepositoryService) EmptyContainerRepository(_ context.Context, opts client.EmptyContainerRepositoryOpts) error {
	err := c.api.EmptyContainerRepository(api.EmptyContainerRepositoryOpts{Name: opts.Name})
	if err != nil {
		return fmt.Errorf("calling API: %w", err)
	}

	return nil
}

// NewContainerRepositoryService returns an initialised container repository service
func NewContainerRepositoryService(
	service api.ContainerRepositoryService,
	state client.ContainerRepositoryState,
	provider v1alpha1.CloudProvider,
) client.ContainerRepositoryService {
	return &containerRepositoryService{
		service:  service,
		state:    state,
		provider: provider,
	}
}
