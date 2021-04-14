package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

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
