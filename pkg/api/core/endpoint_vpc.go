package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateVpcEndpoint(s api.VpcService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateVpc(ctx, request.(api.CreateVpcOpts))
	}
}

func makeDeleteVpcEndpoint(s api.VpcService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Ok(), s.DeleteVpc(ctx, request.(api.DeleteVpcOpts))
	}
}
