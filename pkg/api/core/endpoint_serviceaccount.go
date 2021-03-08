package core // nolint: dupl

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateExternalSecretsServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalSecretsServiceAccount(ctx, request.(api.CreateExternalSecretsServiceAccountOpts))
	}
}

func makeCreateAlbIngressControllerServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAlbIngressControllerServiceAccount(ctx, request.(api.CreateAlbIngressControllerServiceAccountOpts))
	}
}

func makeCreateAWSLoadBalancerControllerServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAWSLoadBalancerControllerServiceAccount(ctx, request.(api.CreateAWSLoadBalancerControllerServiceAccountOpts))
	}
}

func makeCreateExternalDNSServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalDNSServiceAccount(ctx, request.(api.CreateExternalDNSServiceAccountOpts))
	}
}

func makeDeleteExternalSecretsServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteExternalSecretsServiceAccount(ctx, request.(api.ID))
	}
}

func makeDeleteAlbIngressControllerServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteAlbIngressControllerServiceAccount(ctx, request.(api.ID))
	}
}

func makeDeleteAWSLoadBalancerControllerServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteAWSLoadBalancerControllerServiceAccount(ctx, request.(api.ID))
	}
}

func makeDeleteExternalDNSServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteExternalDNSServiceAccount(ctx, request.(api.ID))
	}
}

func makeCreateAutoscalerServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAutoscalerServiceAccount(ctx, request.(api.CreateAutoscalerServiceAccountOpts))
	}
}

func makeDeleteAutoscalerServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteAutoscalerServiceAccount(ctx, request.(api.ID))
	}
}

func makeCreateBlockstorageServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateBlockstorageServiceAccount(ctx, request.(api.CreateBlockstorageServiceAccountOpts))
	}
}

func makeDeleteBlockstorageServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteBlockstorageServiceAccount(ctx, request.(api.ID))
	}
}

func makeCreateServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateServiceAccount(ctx, request.(api.CreateServiceAccountOpts))
	}
}

func makeDeleteServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteServiceAccount(ctx, request.(api.DeleteServiceAccountOpts))
	}
}
