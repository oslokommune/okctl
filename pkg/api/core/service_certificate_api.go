package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type certificateService struct {
	cloudProvider api.CertificateCloudProvider
	store         api.CertificateStore
}

func (c *certificateService) DeleteCertificate(_ context.Context, opts api.DeleteCertificateOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating certificate inputs", errors.Invalid)
	}

	err = c.cloudProvider.DeleteCertificate(opts)
	if err != nil {
		return errors.E(err, "deleting certificate", errors.Internal)
	}

	return nil
}

func (c *certificateService) CreateCertificate(ctx context.Context, opts api.CreateCertificateOpts) (*api.Certificate, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating certificate inputs", errors.Invalid)
	}

	cert, err := c.cloudProvider.CreateCertificate(opts)
	if err != nil {
		return nil, errors.E(err, "creating certificate", errors.Internal)
	}

	err = c.store.SaveCertificate(cert)
	if err != nil {
		return nil, errors.E(err, "storing certificate", errors.IO)
	}

	return cert, nil
}

// NewCertificateService returns an initialised certificate service
func NewCertificateService(cloudProvider api.CertificateCloudProvider, store api.CertificateStore) api.CertificateService {
	return &certificateService{
		cloudProvider: cloudProvider,
		store:         store,
	}
}
