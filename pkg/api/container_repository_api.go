package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CreateContainerRepositoryOpts contains necesessary information to create a container repository
type CreateContainerRepositoryOpts struct {
	ClusterID ID
	Name      string
	StackName string
}

// Validate ensures the struct contains valid data
func (c CreateContainerRepositoryOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClusterID, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.StackName, validation.Required),
	)
}

// DeleteContainerRepositoryOpts contains necessary information to delete a container repository
type DeleteContainerRepositoryOpts struct {
	ClusterID ID
	StackName string
}

// Validate ensures the struct contains valid data
func (c DeleteContainerRepositoryOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClusterID, validation.Required),
		validation.Field(&c.StackName, validation.Required),
	)
}

// EmptyContainerRepositoryOpts contains necessary information to remove all images in a container repository
type EmptyContainerRepositoryOpts struct {
	Name string
}

// Validate ensures the struct contains valid data
func (c EmptyContainerRepositoryOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
	)
}

// ContainerRepository represents a container repository
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
	EmptyContainerRepository(ctx context.Context, opts EmptyContainerRepositoryOpts) error
}

// ContainerRepositoryCloudProvider defines the required cloud operations
// for container repositories
type ContainerRepositoryCloudProvider interface {
	CreateContainerRepository(opts *CreateContainerRepositoryOpts) (*ContainerRepository, error)
	DeleteContainerRepository(opts *DeleteContainerRepositoryOpts) error
	EmptyContainerRepository(ctx context.Context, opts EmptyContainerRepositoryOpts) error
}
