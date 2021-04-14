package core // nolint: dupl

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateServiceAccount(ctx, *request.(*api.CreateServiceAccountOpts))
	}
}

func makeDeleteServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteServiceAccount(ctx, *request.(*api.DeleteServiceAccountOpts))
	}
}
