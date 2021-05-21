package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

type containerRepositoryState struct {
	node breeze.Client
}

// ContainerRepository contains storm compatible state
type ContainerRepository struct {
	Metadata `storm:"inline"`

	ClusterID              ID
	ImageName              string
	StackName              string `storm:"unique"`
	CloudFormationTemplate string
}

// NewContainerRepository returns storm compatible state
func NewContainerRepository(r *client.ContainerRepository, meta Metadata) *ContainerRepository {
	return &ContainerRepository{
		Metadata:               meta,
		ClusterID:              NewID(r.ClusterID),
		ImageName:              r.ImageName,
		StackName:              r.StackName,
		CloudFormationTemplate: r.CloudFormationTemplate,
	}
}

// Convert to client.ContainerRepository
func (r *ContainerRepository) Convert() *client.ContainerRepository {
	return &client.ContainerRepository{
		ClusterID:              r.ClusterID.Convert(),
		ImageName:              r.ImageName,
		StackName:              r.StackName,
		CloudFormationTemplate: r.CloudFormationTemplate,
	}
}

func (c *containerRepositoryState) SaveContainerRepository(repository *client.ContainerRepository) error {
	existing, err := c.getContainerRepository(repository.StackName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return c.node.Save(NewContainerRepository(repository, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return c.node.Save(NewContainerRepository(repository, existing.Metadata))
}

func (c *containerRepositoryState) RemoveContainerRepository(imageName string) error {
	r := &ContainerRepository{}

	err := c.node.One("ImageName", imageName, r)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return c.node.DeleteStruct(r)
}

func (c *containerRepositoryState) GetContainerRepository(imageName string) (*client.ContainerRepository, error) {
	r, err := c.getContainerRepository(imageName)
	if err != nil {
		return nil, err
	}

	return r.Convert(), nil
}

func (c *containerRepositoryState) getContainerRepository(imageName string) (*ContainerRepository, error) {
	r := &ContainerRepository{}

	err := c.node.One("ImageName", imageName, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NewContainerRepositoryState returns an initialised state client
func NewContainerRepositoryState(node breeze.Client) client.ContainerRepositoryState {
	return &containerRepositoryState{
		node: node,
	}
}
