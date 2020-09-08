package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type certificateService struct {
	api   client.CertificateAPI
	store client.CertificateStore
}

func (s *certificateService) CreateCertificate(_ context.Context, opts api.CreateCertificateOpts) (*api.Certificate, error) {
	certificate, err := s.api.CreateCertificate(opts)
	if err != nil {
		return nil, err
	}

	_, err = s.store.SaveCertificate(certificate)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

// NewCertificateService returns an initialised service
func NewCertificateService(api client.CertificateAPI, store client.CertificateStore) client.CertificateService {
	return &certificateService{
		api:   api,
		store: store,
	}
}
