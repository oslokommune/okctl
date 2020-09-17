package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateExternalSecretsPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalSecretsPolicy(ctx, request.(api.CreateExternalSecretsPolicyOpts))
	}
}

func makeCreateAlbIngressControllerPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAlbIngressControllerPolicy(ctx, request.(api.CreateAlbIngressControllerPolicyOpts))
	}
}

func makeCreateExternalDNSPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalDNSPolicy(ctx, request.(api.CreateExternalDNSPolicyOpts))
	}
}

func makeDeleteExternalSecretsPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteExternalSecretsPolicy(ctx, request.(api.ID))
	}
}

func makeDeleteAlbIngressControllerPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteAlbIngressControllerPolicy(ctx, request.(api.ID))
	}
}

func makeDeleteExternalDNSPolicyEndpoint(s api.ManagedPolicyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteExternalDNSPolicy(ctx, request.(api.ID))
	}
}
