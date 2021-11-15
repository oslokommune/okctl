package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateClusterEndpoint(s api.ClusterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateCluster(ctx, request.(api.ClusterCreateOpts))
	}
}

func makeDeleteClusterEndpoint(s api.ClusterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Ok(), s.DeleteCluster(ctx, request.(api.ClusterDeleteOpts))
	}
}

func makeGetClusterSecurityGroupIDEndpoint(s api.ClusterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetClusterSecurityGroupID(ctx, request.(*api.ClusterSecurityGroupIDGetOpts))
	}
}
