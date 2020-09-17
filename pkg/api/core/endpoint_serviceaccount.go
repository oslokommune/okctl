package core

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

func makeDeleteExternalDNSServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteExternalDNSServiceAccount(ctx, request.(api.ID))
	}
}
