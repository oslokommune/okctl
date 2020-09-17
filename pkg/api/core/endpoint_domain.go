package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateHostedZoneEndpoint(s api.DomainService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateHostedZone(ctx, request.(api.CreateHostedZoneOpts))
	}
}

func makeDeleteHostedZoneEndpoint(s api.DomainService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteHostedZone(ctx, request.(api.DeleteHostedZoneOpts))
	}
}
