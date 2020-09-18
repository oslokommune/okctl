package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type certificateService struct {
	spinner spinner.Spinner
	api     client.CertificateAPI
	store   client.CertificateStore
	state   client.CertificateState
	report  client.CertificateReport
}

func (s *certificateService) CreateCertificate(_ context.Context, opts api.CreateCertificateOpts) (*api.Certificate, error) {
	err := s.spinner.Start("certificate")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	c := s.state.GetCertificate(opts.Domain)
	if c.Validate() == nil {
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
	spinner spinner.Spinner,
	api client.CertificateAPI,
	store client.CertificateStore,
	state client.CertificateState,
	report client.CertificateReport,
) client.CertificateService {
	return &certificateService{
		spinner: spinner,
		api:     api,
		store:   store,
		state:   state,
		report:  report,
	}
}
