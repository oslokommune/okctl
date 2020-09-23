package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateIdentityPoolEndpoint(s api.IdentityManagerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateIdentityPool(ctx, request.(api.CreateIdentityPoolOpts))
	}
}
