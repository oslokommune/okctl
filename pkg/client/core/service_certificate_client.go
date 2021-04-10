package core

import (
	"context"
	"errors"

	stormpkg "github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type certificateService struct {
	api   client.CertificateAPI
	state client.CertificateState
}

func (s *certificateService) DeleteCognitoCertificate(_ context.Context, opts client.DeleteCognitoCertificateOpts) error {
	err := s.api.DeleteCognitoCertificate(api.DeleteCognitoCertificateOpts{
		ID:     opts.ID,
		Domain: opts.Domain,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *certificateService) DeleteCertificate(_ context.Context, opts client.DeleteCertificateOpts) error {
	err := s.api.DeleteCertificate(api.DeleteCertificateOpts{
		ID:     opts.ID,
		Domain: opts.Domain,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveCertificate(opts.Domain)
	if err != nil {
		return err
	}

	return nil
}

func (s *certificateService) CreateCertificate(_ context.Context, opts client.CreateCertificateOpts) (*client.Certificate, error) {
	// [Refactor] Reconciler is responsible for ordering operations
	//
	// We should be doing this check in the reconciler together with a
	// verification towards the AWS API. Keeping this here for the
	// time being, so we are compatible with expected behavior.
	{
		c, err := s.state.GetCertificate(opts.Domain)
		if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
			return nil, err
		}

		if err == nil {
			return c, nil
		}
	}

	c, err := s.api.CreateCertificate(api.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         opts.FQDN,
		Domain:       opts.Domain,
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	certificate := &client.Certificate{
		ID:                     c.ID,
		FQDN:                   c.FQDN,
		Domain:                 c.Domain,
		HostedZoneID:           c.HostedZoneID,
		ARN:                    c.CertificateARN,
		StackName:              c.StackName,
		CloudFormationTemplate: c.CloudFormationTemplate,
	}

	err = s.state.SaveCertificate(certificate)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

// NewCertificateService returns an initialised service
func NewCertificateService(
	api client.CertificateAPI,
	state client.CertificateState,
) client.CertificateService {
	return &certificateService{
		api:   api,
		state: state,
	}
}
