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
	LambdaPolicyARN              string
	LambdaRoleARN                string
	LambdaFunctionARN            string
	CloudFormationTemplate       string
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
	RotaterBucket     string
	RotaterKey        string
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

// S3Bucket contains the state after an AWS S3
// bucket has been created
type S3Bucket struct {
	ID                     ID
	Name                   string
	StackName              string
	CloudFormationTemplate string
}

// CreateS3BucketOpts contains the required inputs
type CreateS3BucketOpts struct {
	ID        ID
	Name      string
	StackName string
}

// Validate the inputs
func (o CreateS3BucketOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.StackName, validation.Required),
	)
}

// DeleteS3BucketOpts contains the required inputs
type DeleteS3BucketOpts struct {
	ID        ID
	StackName string
}

// Validate the inputs
func (o DeleteS3BucketOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.StackName, validation.Required),
	)
}

// ComponentService defines a set of operations for creating
// components that integrate with the Kubernetes cluster's
// applications
type ComponentService interface {
	CreatePostgresDatabase(ctx context.Context, opts *CreatePostgresDatabaseOpts) (*PostgresDatabase, error)
	DeletePostgresDatabase(ctx context.Context, opts *DeletePostgresDatabaseOpts) error
	CreateS3Bucket(ctx context.Context, opts *CreateS3BucketOpts) (*S3Bucket, error)
	DeleteS3Bucket(ctx context.Context, opts *DeleteS3BucketOpts) error
}

// ComponentCloudProvider defines the required cloud operations
// for the components
type ComponentCloudProvider interface {
	CreatePostgresDatabase(ctx context.Context, opts *CreatePostgresDatabaseOpts) (*PostgresDatabase, error)
	DeletePostgresDatabase(opts *DeletePostgresDatabaseOpts) error
	CreateS3Bucket(ctx context.Context, opts *CreateS3BucketOpts) (*S3Bucket, error)
	DeleteS3Bucket(opts *DeleteS3BucketOpts) error
}
