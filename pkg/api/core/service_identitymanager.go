package core

import (
	"context"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/api"
)

type identityManagerService struct {
	provider api.IdentityManagerCloudProvider
	cert     api.CertificateCloudProvider
}

func (s *identityManagerService) CreateIdentityPoolClient(_ context.Context, opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	client, err := s.provider.CreateIdentityPoolClient(opts)
	if err != nil {
		return nil, errors.E(err, "creating an identity pool client: %w", errors.Internal)
	}

	return client, nil
}

func (s *identityManagerService) CreateIdentityPool(_ context.Context, opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	certificate, err := s.cert.CreateCertificate(api.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         opts.AuthFQDN,
		Domain:       opts.AuthDomain,
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, errors.E(err, "creating a certificate for auth domain", errors.Internal)
	}

	pool, err := s.provider.CreateIdentityPool(certificate.CertificateARN, opts)
	if err != nil {
		return nil, errors.E(err, "creating an identity pool", errors.Internal)
	}

	pool.Certificate = certificate

	return pool, nil
}

// NewIdentityManagerService returns an initialised identity manager
func NewIdentityManagerService(
	provider api.IdentityManagerCloudProvider,
	cert api.CertificateCloudProvider,
) api.IdentityManagerService {
	return &identityManagerService{
		provider: provider,
		cert:     cert,
	}
}
