package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type componentDirectClient struct {
	service api.ComponentService
}

func (c *componentDirectClient) CreateS3Bucket(opts api.CreateS3BucketOpts) (*api.S3Bucket, error) {
	return c.service.CreateS3Bucket(context.Background(), &opts)
}

func (c *componentDirectClient) DeleteS3Bucket(opts api.DeleteS3BucketOpts) error {
	return c.service.DeleteS3Bucket(context.Background(), &opts)
}

func (c *componentDirectClient) CreatePostgresDatabase(opts api.CreatePostgresDatabaseOpts) (*api.PostgresDatabase, error) {
	return c.service.CreatePostgresDatabase(context.Background(), &opts)
}

func (c *componentDirectClient) DeletePostgresDatabase(opts api.DeletePostgresDatabaseOpts) error {
	return c.service.DeletePostgresDatabase(context.Background(), &opts)
}

// NewComponentAPI returns an initialised REST API invoker
func NewComponentAPI(service api.ComponentService) client.ComponentAPI {
	return &componentDirectClient{
		service: service,
	}
}
