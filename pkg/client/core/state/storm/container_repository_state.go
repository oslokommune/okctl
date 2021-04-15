package storm

import (
	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type containerRepositoryState struct {
	node stormpkg.Node
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
	return c.node.Save(NewContainerRepository(repository, NewMetadata()))
}

func (c *containerRepositoryState) RemoveContainerRepository(stackName string) error {
	r := &ContainerRepository{}

	err := c.node.One("StackName", stackName, r)
	if err != nil {
		return err
	}

	return c.node.DeleteStruct(r)
}

func (c *containerRepositoryState) GetContainerRepository(stackName string) (*client.ContainerRepository, error) {
	r := &ContainerRepository{}

	err := c.node.One("StackName", stackName, r)
	if err != nil {
		return nil, err
	}

	return r.Convert(), nil
}

// NewContainerRepositoryState returns an initialised state client
func NewContainerRepositoryState(node stormpkg.Node) client.ContainerRepositoryState {
	return &containerRepositoryState{
		node: node,
	}
}
