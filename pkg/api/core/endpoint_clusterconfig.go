package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateClusterConfigEndpoint(s api.ClusterConfigService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateClusterConfig(ctx, request.(api.CreateClusterConfigOpts))
	}
}
