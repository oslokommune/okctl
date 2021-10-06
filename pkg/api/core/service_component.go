package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type componentService struct {
	provider api.ComponentCloudProvider
}

func (c *componentService) CreateS3Bucket(ctx context.Context, opts *api.CreateS3BucketOpts) (*api.S3Bucket, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "invalid inputs", errors.Invalid)
	}

	b, err := c.provider.CreateS3Bucket(ctx, opts)
	if err != nil {
		return nil, errors.E(err, "creating S3 bucket", errors.Internal)
	}

	return b, nil
}

func (c *componentService) DeleteS3Bucket(_ context.Context, opts *api.DeleteS3BucketOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "invalid inputs", errors.Invalid)
	}

	err = c.provider.DeleteS3Bucket(opts)
	if err != nil {
		return errors.E(err, "deleting S3 bucket", errors.Internal)
	}

	return nil
}

func (c *componentService) CreatePostgresDatabase(ctx context.Context, opts *api.CreatePostgresDatabaseOpts) (*api.PostgresDatabase, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "invalid inputs", errors.Invalid)
	}

	pg, err := c.provider.CreatePostgresDatabase(ctx, opts)
	if err != nil {
		return nil, errors.E(err, "creating postgres database", errors.Internal)
	}

	return pg, nil
}

func (c *componentService) DeletePostgresDatabase(_ context.Context, opts *api.DeletePostgresDatabaseOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "invalid inputs", errors.Invalid)
	}

	err = c.provider.DeletePostgresDatabase(opts)
	if err != nil {
		return errors.E(err, "deleting postgres database", errors.Internal)
	}

	return nil
}

// NewComponentService returns an initialised component service
func NewComponentService(provider api.ComponentCloudProvider) api.ComponentService {
	return &componentService{
		provider: provider,
	}
}
