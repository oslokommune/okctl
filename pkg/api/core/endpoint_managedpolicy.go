package core // nolint: dupl

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateExternalDNSPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalDNSPolicy(ctx, request.(api.CreateExternalDNSPolicyOpts))
	}
}

func makeDeleteExternalDNSPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteExternalDNSPolicy(ctx, request.(api.ID))
	}
}

func makeCreateAutoscalerPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAutoscalerPolicy(ctx, request.(api.CreateAutoscalerPolicy))
	}
}

func makeDeleteAutoscalerPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteAutoscalerPolicy(ctx, request.(api.ID))
	}
}

func makeCreateBlockstoragePolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateBlockstoragePolicy(ctx, request.(api.CreateBlockstoragePolicy))
	}
}

func makeDeleteBlockstoragePolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteBlockstoragePolicy(ctx, request.(api.ID))
	}
}

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
