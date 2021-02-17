package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateExternalSecretsHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalSecretsHelmChart(ctx, request.(api.CreateExternalSecretsHelmChartOpts))
	}
}

func makeCreateKubePrometheusStackHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateKubePrometheusStack(ctx, request.(api.CreateKubePrometheusStackOpts))
	}
}

func makeCreateAlbIngressControllerHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAlbIngressControllerHelmChart(ctx, request.(api.CreateAlbIngressControllerHelmChartOpts))
	}
}

func makeCreateAWSLoadBalancerControllerHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAWSLoadBalancerControllerHelmChart(ctx, request.(api.CreateAWSLoadBalancerControllerHelmChartOpts))
	}
}

func makeCreateArgoCD(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateArgoCD(ctx, request.(api.CreateArgoCDOpts))
	}
}

func makeCreateAutoscalerHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAutoscalerHelmChart(ctx, request.(api.CreateAutoscalerHelmChartOpts))
	}
}

func makeCreateBlockstorageHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateBlockstorageHelmChart(ctx, request.(api.CreateBlockstorageHelmChartOpts))
	}
}

func makeCreateKubePrometheusStack(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateKubePrometheusStack(ctx, request.(api.CreateKubePrometheusStackOpts))
	}
}
