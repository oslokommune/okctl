package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateExternalSecretsServiceAccountEndpoint(s api.ServiceAccountService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateExternalSecretsServiceAccount(ctx, request.(api.CreateExternalSecretsServiceAccountOpts))
	}
}
