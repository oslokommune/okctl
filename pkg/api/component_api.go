package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// PostgresDatabase contains the state resulting
// from a newly created Postgres database
type PostgresDatabase struct {
	ID                           ID
	ApplicationName              string
	UserName                     string
	StackName                    string
	AdminSecretFriendlyName      string
	EndpointAddress              string
	EndpointPort                 int
	OutgoingSecurityGroupID      string
	SecretsManagerAdminSecretARN string
	CloudFormationTemplate       []byte
}

// CreatePostgresDatabaseOpts contains the inputs
// required for creating a Postgres database
type CreatePostgresDatabaseOpts struct {
	ID                ID
	ApplicationName   string
	UserName          string
	StackName         string
	VpcID             string
	DBSubnetGroupName string
	DBSubnetIDs       []string
	DBSubnetCIDRs     []string
}

// Validate the inputs
func (o *CreatePostgresDatabaseOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID),
	)
}

// DeletePostgresDatabaseOpts contains the inputs
// required for deleting a Postgres database
type DeletePostgresDatabaseOpts struct {
	ID        ID
	StackName string
}

// Validate the inputs
func (o *DeletePostgresDatabaseOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID),
		validation.Field(&o.StackName),
	)
}

// ComponentService defines a set of operations for creating
// components that integrate with the Kubernetes cluster's
// applications
type ComponentService interface {
	CreatePostgresDatabase(ctx context.Context, opts *CreatePostgresDatabaseOpts) (*PostgresDatabase, error)
	DeletePostgresDatabase(ctx context.Context, opts *DeletePostgresDatabaseOpts) error
}

// ComponentCloudProvider defines the required cloud operations
// for the components
type ComponentCloudProvider interface {
	CreatePostgresDatabase(opts *CreatePostgresDatabaseOpts) (*PostgresDatabase, error)
	DeletePostgresDatabase(opts *DeletePostgresDatabaseOpts) error
}
