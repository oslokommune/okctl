package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateIdentityPoolEndpoint(s api.IdentityManagerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateIdentityPool(ctx, request.(api.CreateIdentityPoolOpts))
	}
}

func makeCreateIdentityPoolClient(s api.IdentityManagerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateIdentityPoolClient(ctx, request.(api.CreateIdentityPoolClientOpts))
	}
}

func makeCreateIdentityPoolUser(s api.IdentityManagerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateIdentityPoolUser(ctx, request.(api.CreateIdentityPoolUserOpts))
	}
}

func makeDeleteIdentityPoolUserEndpoint(s api.IdentityManagerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteIdentityPoolUser(ctx, request.(api.DeleteIdentityPoolUserOpts))
	}
}

func makeDeleteIdentityPoolEndpoint(s api.IdentityManagerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteIdentityPool(ctx, request.(api.DeleteIdentityPoolOpts))
	}
}

func makeDeleteIdentityPoolClientEndpoint(s api.IdentityManagerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteIdentityPoolClient(ctx, request.(api.DeleteIdentityPoolClientOpts))
	}
}
