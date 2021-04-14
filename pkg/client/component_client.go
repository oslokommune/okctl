package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// CreatePostgresDatabaseOpts contains the required inputs
type CreatePostgresDatabaseOpts struct {
	ID                api.ID
	ApplicationName   string
	UserName          string
	VpcID             string
	DBSubnetGroupName string
	DBSubnetIDs       []string
	DBSubnetCIDRs     []string
	Namespace         string
}

// DeletePostgresDatabaseOpts contains the required inputs
type DeletePostgresDatabaseOpts struct {
	ID              api.ID
	ApplicationName string
	VpcID           string
}

// PostgresDatabase contains the state after
// creating the database
type PostgresDatabase struct {
	ID                           api.ID
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
	Namespace                    string
	AdminSecretName              string
	AdminSecretARN               string
	DatabaseConfigMapName        string
	RotaterBucket                *S3Bucket
}

// S3Bucket contains the state after an AWS S3
// bucket has been created
type S3Bucket struct {
	Name                   string
	StackName              string
	CloudFormationTemplate string
}

// ComponentService orchestrates the creation of various services
type ComponentService interface {
	CreatePostgresDatabase(ctx context.Context, opts CreatePostgresDatabaseOpts) (*PostgresDatabase, error)
	DeletePostgresDatabase(ctx context.Context, opts DeletePostgresDatabaseOpts) error
}

// ComponentAPI invokes the API
type ComponentAPI interface {
	CreatePostgresDatabase(opts api.CreatePostgresDatabaseOpts) (*api.PostgresDatabase, error)
	DeletePostgresDatabase(opts api.DeletePostgresDatabaseOpts) error
	CreateS3Bucket(opts api.CreateS3BucketOpts) (*api.S3Bucket, error)
	DeleteS3Bucket(opts api.DeleteS3BucketOpts) error
}

// ComponentState updates the state
type ComponentState interface {
	SavePostgresDatabase(database *PostgresDatabase) error
	RemovePostgresDatabase(stackName string) error
	GetPostgresDatabase(stackName string) (*PostgresDatabase, error)
}
