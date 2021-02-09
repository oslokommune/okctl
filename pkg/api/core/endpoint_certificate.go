package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateCertificateEndpoint(s api.CertificateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.CreateCertificate(ctx, request.(api.CreateCertificateOpts))
	}
}

func makeDeleteCertificateEndpoint(s api.CertificateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Empty{}, s.DeleteCertificate(ctx, request.(api.DeleteCertificateOpts))
	}
}
