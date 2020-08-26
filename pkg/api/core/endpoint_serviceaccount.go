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

func makeCreateExternalDnsServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalDnsServiceAccount(ctx, request.(api.CreateExternalDnsServiceAccountOpts))
	}
}
