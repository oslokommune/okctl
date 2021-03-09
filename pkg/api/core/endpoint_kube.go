package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateExternalDNSKubeDeploymentEndpoint(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalDNSKubeDeployment(ctx, request.(api.CreateExternalDNSKubeDeploymentOpts))
	}
}

func makeDeleteNamespaceEndpoint(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteNamespace(ctx, request.(api.DeleteNamespaceOpts))
	}
}

func makeCreateStorageClass(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateStorageClass(ctx, request.(api.CreateStorageClassOpts))
	}
}

func makeCreateExternalSecretsEndpoint(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalSecrets(ctx, request.(api.CreateExternalSecretsOpts))
	}
}

func makeDeleteExternalSecrets(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteExternalSecrets(ctx, request.(api.DeleteExternalSecretsOpts))
	}
}

func makeCreateConfigMapEndpoint(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateConfigMap(ctx, request.(api.CreateConfigMapOpts))
	}
}

func makeDeleteConfigMap(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteConfigMap(ctx, request.(api.DeleteConfigMapOpts))
	}
}

func makeScaleDeployment(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.ScaleDeployment(ctx, request.(api.ScaleDeploymentOpts))
	}
}

func makeCreateNamespace(s api.KubeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateNamespace(ctx, *request.(*api.CreateNamespaceOpts))
	}
}
