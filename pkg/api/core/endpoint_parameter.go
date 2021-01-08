package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateSecret(s api.ParameterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateSecret(ctx, request.(api.CreateSecretOpts))
	}
}

func makeDeleteSecret(s api.ParameterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return &Empty{}, s.DeleteSecret(ctx, request.(api.DeleteSecretOpts))
	}
}
