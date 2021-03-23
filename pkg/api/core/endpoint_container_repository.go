package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateContainerRepositoryEndpoint(s api.ContainerRepositoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateContainerRepository(ctx, request.(*api.CreateContainerRepositoryOpts))
	}
}

func makeDeleteContainerRepositoryEndpoint(s api.ContainerRepositoryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteContainerRepository(ctx, request.(*api.DeleteContainerRepositoryOpts))
	}
}
