package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type containerRepositoryService struct {
	spinner spinner.Spinner
	api     client.ContainerRepositoryAPI
	store   client.ContainerRepositoryStore
	state   client.ContainerRepositoryState
	report  client.ContainerRepositoryReport

	provider v1alpha1.CloudProvider
}

// CreateContainerRepository handles api, state, store and report orchistration for creation of a container repository
func (c *containerRepositoryService) CreateContainerRepository(_ context.Context, opts client.CreateContainerRepositoryOpts) (*client.ContainerRepository, error) {
	err := c.spinner.Start("container registry")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	repository, err := c.api.CreateContainerRepository(api.CreateContainerRepositoryOpts{
		ClusterID: opts.ClusterID,
		Name:      opts.ImageName,
		StackName: cfn.NewStackNamer().ContainerRepository(opts.ImageName, opts.ClusterID.Repository, opts.ClusterID.Environment),
	})
	if err != nil {
		return nil, err
	}

	containerRepository := &client.ContainerRepository{
		ImageName: repository.Name,
	}

	r1, err := c.store.SaveContainerRepository(containerRepository)
	if err != nil {
		return nil, err
	}

	r2, err := c.state.SaveContainerRepository(containerRepository)
	if err != nil {
		return nil, err
	}

	err = c.report.ReportCreateContainerRepository(containerRepository, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return containerRepository, nil
}

// DeleteContainerRepository handles api, state, store and report orchistration for deletion of a container repository
func (c *containerRepositoryService) DeleteContainerRepository(_ context.Context, opts client.DeleteContainerRepositoryOpts) error {
	err := c.spinner.Start("container registry")
	if err != nil {
		return err
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	err = c.api.DeleteContainerRepository(api.DeleteContainerRepositoryOpts{
		ClusterID: opts.ClusterID,
		StackName: cfn.NewStackNamer().ContainerRepository(opts.ImageName, opts.ClusterID.Repository, opts.ClusterID.Environment),
	})
	if err != nil {
		return err
	}

	r1, err := c.store.RemoveContainerRepository(opts.ImageName)
	if err != nil {
		return err
	}

	r2, err := c.state.RemoveContainerRepository(opts.ImageName)
	if err != nil {
		return err
	}

	err = c.report.ReportDeleteContainerRepository(opts.ImageName, []*store.Report{r1, r2})
	if err != nil {
		return err
	}

	return nil
}

// NewContainerRepositoryService returns an initialised container repository service
func NewContainerRepositoryService(
	spin spinner.Spinner,
	api client.ContainerRepositoryAPI,
	store client.ContainerRepositoryStore,
	state client.ContainerRepositoryState,
	report client.ContainerRepositoryReport,
	provider v1alpha1.CloudProvider,
) client.ContainerRepositoryService {
	return &containerRepositoryService{
		spinner:  spin,
		api:      api,
		store:    store,
		state:    state,
		report:   report,
		provider: provider,
	}
}
