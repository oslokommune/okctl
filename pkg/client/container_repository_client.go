package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type CreateContainerRepositoryOpts struct {
	ClusterID api.ID
	ImageName string
}

type DeleteContainerRepositoryOpts struct {
	ClusterID api.ID
	ImageName string
}

// ContainerRepository contains state after creating a container repository
type ContainerRepository struct {
	api.ContainerRepository
	ImageName string
}

// ContainerRepositoryService orchestrates the creation of various services
type ContainerRepositoryService interface {
	CreateContainerRepository(ctx context.Context, opts CreateContainerRepositoryOpts) (*ContainerRepository, error)
	DeleteContainerRepository(ctx context.Context, opts DeleteContainerRepositoryOpts) error
}

// ContainerRepositoryAPI invokes the API
type ContainerRepositoryAPI interface {
	CreateContainerRepository(opts api.CreateContainerRepositoryOpts) (*api.ContainerRepository, error)
	DeleteContainerRepository(opts api.DeleteContainerRepositoryOpts) error
}

// ContainerRepositoryStore saves the data
type ContainerRepositoryStore interface {
	SaveContainerRepository(repository *ContainerRepository) (*store.Report, error)
	RemoveContainerRepository(imageName string) (*store.Report, error)
}

// ContainerRepositoryState updates the state
type ContainerRepositoryState interface {
	SaveContainerRepository(repository *ContainerRepository) (*store.Report, error)
	RemoveContainerRepository(imageName string) (*store.Report, error)
	GetContainerRepository(imageName string) (*ContainerRepository, error)
}

// ContainerRepositoryReport reports on the state and storage operations
type ContainerRepositoryReport interface {
	ReportCreateContainerRepository(repository *ContainerRepository, reports []*store.Report) error
	ReportDeleteContainerRepository(imageName string, reports []*store.Report) error
}
