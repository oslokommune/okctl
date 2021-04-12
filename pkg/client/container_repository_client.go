package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// CreateContainerRepositoryOpts contains necessary information to create a container repository
type CreateContainerRepositoryOpts struct {
	ClusterID api.ID
	ImageName string
}

// DeleteContainerRepositoryOpts contains necessary information to delete a container repository
type DeleteContainerRepositoryOpts struct {
	ClusterID api.ID
	ImageName string
}

// ContainerRepository contains state after creating a container repository
type ContainerRepository struct {
	ClusterID              api.ID
	ImageName              string
	StackName              string
	CloudFormationTemplate string
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

// ContainerRepositoryState updates the state
type ContainerRepositoryState interface {
	SaveContainerRepository(repository *ContainerRepository) (*store.Report, error)
	RemoveContainerRepository(imageName string) (*store.Report, error)
	GetContainerRepository(imageName string) (*ContainerRepository, error)
}

// URI returns the URI where the image can be pulled and pushed
func (c ContainerRepository) URI() url.URL {
	return url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", c.ClusterID.AWSAccountID, c.ClusterID.Region),
		Path:   fmt.Sprintf("%s-%s", c.ClusterID.ClusterName, c.ImageName),
	}
}
