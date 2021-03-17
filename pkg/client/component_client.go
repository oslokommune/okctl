package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
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
	Namespace             string
	AdminSecretName       string
	DatabaseConfigMapName string
	RotaterBucket         *api.S3Bucket
	*api.PostgresDatabase
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

// ComponentStore saves the data
type ComponentStore interface {
	SavePostgresDatabase(database *PostgresDatabase) (*store.Report, error)
	RemovePostgresDatabase(applicationName string) (*store.Report, error)
}

// ComponentState updates the state
type ComponentState interface {
	SavePostgresDatabase(database *PostgresDatabase) (*store.Report, error)
	RemovePostgresDatabase(applicationName string) (*store.Report, error)
	GetPostgresDatabase(applicationName string) (*PostgresDatabase, error)
}

// ComponentReport reports on the state and storage operations
type ComponentReport interface {
	ReportCreatePostgresDatabase(database *PostgresDatabase, reports []*store.Report) error
	ReportDeletePostgresDatabase(applicationName string, reports []*store.Report) error
}
