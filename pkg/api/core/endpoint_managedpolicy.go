package core // nolint: dupl

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreatePolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreatePolicy(ctx, *request.(*api.CreatePolicyOpts))
	}
}

func makeDeletePolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeletePolicy(ctx, *request.(*api.DeletePolicyOpts))
	}
}
