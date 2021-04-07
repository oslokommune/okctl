package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateArgoCD(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateArgoCD(ctx, request.(api.CreateArgoCDOpts))
	}
}

func makeCreateKubePrometheusStack(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateKubePrometheusStack(ctx, request.(api.CreateKubePrometheusStackOpts))
	}
}

func makeCreateLokiHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateLokiHelmChart(ctx, request.(api.CreateLokiHelmChartOpts))
	}
}

func makeCreatePromtailHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreatePromtailHelmChart(ctx, request.(api.CreatePromtailHelmChartOpts))
	}
}

func makeCreateHelmRelease(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateHelmRelease(ctx, request.(api.CreateHelmReleaseOpts))
	}
}

func makeDeleteHelmRelease(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteHelmRelease(ctx, request.(api.DeleteHelmReleaseOpts))
	}
}
