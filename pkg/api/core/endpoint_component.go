package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreatePostgresDatabaseEndpoint(s api.ComponentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreatePostgresDatabase(ctx, request.(*api.CreatePostgresDatabaseOpts))
	}
}

func makeDeletePostgresDatabaseEndpoint(s api.ComponentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeletePostgresDatabase(ctx, request.(*api.DeletePostgresDatabaseOpts))
	}
}

func makeCreateS3BucketEndpoint(s api.ComponentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateS3Bucket(ctx, request.(*api.CreateS3BucketOpts))
	}
}

func makeDeleteS3BucketEndpoint(s api.ComponentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteS3Bucket(ctx, request.(*api.DeleteS3BucketOpts))
	}
}
