package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type certificateService struct {
	api    client.CertificateAPI
	store  client.CertificateStore
	state  client.CertificateState
	report client.CertificateReport
}

func (s *certificateService) CreateCertificate(_ context.Context, opts api.CreateCertificateOpts) (*api.Certificate, error) {
	c := s.state.GetCertificate(opts.Domain)
	if len(c.CertificateARN) > 0 {
		// We should already have this certificate, if not
		// something is very wrong.
		return s.store.GetCertificate(opts.Domain)
	}

	certificate, err := s.api.CreateCertificate(opts)
	if err != nil {
		return nil, err
	}

	r1, err := s.store.SaveCertificate(certificate)
	if err != nil {
		return nil, err
	}

	r2, err := s.state.SaveCertificate(certificate)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveCertificate(certificate, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

// NewCertificateService returns an initialised service
func NewCertificateService(
	api client.CertificateAPI,
	store client.CertificateStore,
	state client.CertificateState,
	report client.CertificateReport,
) client.CertificateService {
	return &certificateService{
		api:    api,
		store:  store,
		state:  state,
		report: report,
	}
}
