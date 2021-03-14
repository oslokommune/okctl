package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetComponentPostgres matches the REST API route
	TargetComponentPostgres = "components/postgres/"
	// TargetComponentS3Bucket matches the REST API route
	TargetComponentS3Bucket = "components/s3bucket/"
)

type componentAPI struct {
	client *HTTPClient
}

func (c *componentAPI) CreateS3Bucket(opts api.CreateS3BucketOpts) (*api.S3Bucket, error) {
	into := &api.S3Bucket{}
	return into, c.client.DoPost(TargetComponentS3Bucket, &opts, into)
}

func (c *componentAPI) DeleteS3Bucket(opts api.DeleteS3BucketOpts) error {
	return c.client.DoDelete(TargetComponentS3Bucket, &opts)
}

func (c *componentAPI) CreatePostgresDatabase(opts api.CreatePostgresDatabaseOpts) (*api.PostgresDatabase, error) {
	into := &api.PostgresDatabase{}
	return into, c.client.DoPost(TargetComponentPostgres, &opts, into)
}

func (c *componentAPI) DeletePostgresDatabase(opts api.DeletePostgresDatabaseOpts) error {
	return c.client.DoDelete(TargetComponentPostgres, &opts)
}

// NewComponentAPI returns an initialised REST API invoker
func NewComponentAPI(client *HTTPClient) client.ComponentAPI {
	return &componentAPI{
		client: client,
	}
}
