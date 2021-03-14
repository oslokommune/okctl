package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type certificateService struct {
	cloudProvider api.CertificateCloudProvider
}

func (c *certificateService) DeleteCognitoCertificate(_ context.Context, opts api.DeleteCognitoCertificateOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = c.cloudProvider.DeleteCognitoCertificate(opts)
	if err != nil {
		return errors.E(err, "deleting cognito certificate", errors.Internal)
	}

	return nil
}

func (c *certificateService) DeleteCertificate(_ context.Context, opts api.DeleteCertificateOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = c.cloudProvider.DeleteCertificate(opts)
	if err != nil {
		return errors.E(err, "deleting certificate", errors.Internal)
	}

	return nil
}

func (c *certificateService) CreateCertificate(_ context.Context, opts api.CreateCertificateOpts) (*api.Certificate, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	cert, err := c.cloudProvider.CreateCertificate(opts)
	if err != nil {
		return nil, errors.E(err, "creating certificate", errors.Internal)
	}

	return cert, nil
}

// NewCertificateService returns an initialised certificate service
func NewCertificateService(cloudProvider api.CertificateCloudProvider) api.CertificateService {
	return &certificateService{
		cloudProvider: cloudProvider,
	}
}
