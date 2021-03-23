package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateContainerRepositoryOpts struct {
	ClusterID ID
	Name      string
	StackName string
}

func (c CreateContainerRepositoryOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClusterID, validation.Required),
		validation.Field(&c.StackName, validation.Required),
	)
}

type DeleteContainerRepositoryOpts struct {
	ClusterID ID
	StackName string
}

func (c DeleteContainerRepositoryOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClusterID, validation.Required),
		validation.Field(&c.StackName, validation.Required),
	)
}

type ContainerRepository struct {
	ClusterID              ID
	Name                   string
	StackName              string
	CloudFormationTemplate string
}

// ContainerRepositoryService defines operations for container repositories
type ContainerRepositoryService interface {
	CreateContainerRepository(ctx context.Context, opts *CreateContainerRepositoryOpts) (*ContainerRepository, error)
	DeleteContainerRepository(ctx context.Context, opts *DeleteContainerRepositoryOpts) error
}

// ContainerRepositoryCloudProvider defines the required cloud operations
// for container repositories
type ContainerRepositoryCloudProvider interface {
	CreateContainerRepository(opts *CreateContainerRepositoryOpts) (*ContainerRepository, error)
	DeleteContainerRepository(opts *DeleteContainerRepositoryOpts) error
}
