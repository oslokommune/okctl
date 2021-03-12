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
