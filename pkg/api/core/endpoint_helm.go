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

func makeCreateAlbIngressControllerHelmChartEndpoint(s api.HelmService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateAlbIngressControllerHelmChart(ctx, request.(api.CreateAlbIngressControllerHelmChartOpts))
	}
}
