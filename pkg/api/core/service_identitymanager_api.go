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

func (s *identityManagerService) DeleteIdentityPool(_ context.Context, opts api.DeleteIdentityPoolOpts) error {
	err := s.provider.DeleteIdentityPool(opts)
	if err != nil {
		return errors.E(err, "deleting an identity pool", errors.Internal)
	}

	err = s.cert.DeleteCertificate(api.DeleteCertificateOpts{
		Domain: opts.Domain,
	})
	if err != nil {
		return errors.E(err, "deleting an identity pool certificate", errors.Internal)
	}

	return nil
}

func (s *identityManagerService) DeleteIdentityPoolClient(_ context.Context, opts api.DeleteIdentityPoolClientOpts) error {
	err := s.provider.DeleteIdentityPoolClient(opts)
	if err != nil {
		return errors.E(err, "deleting identity pool client", errors.Internal)
	}

	return nil
}

func (s *identityManagerService) CreateIdentityPoolUser(_ context.Context, opts api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error) {
	user, err := s.provider.CreateIdentityPoolUser(opts)
	if err != nil {
		return nil, errors.E(err, "creating an identity pool user", errors.Internal)
	}

	return user, nil
}

func (s *identityManagerService) CreateIdentityPoolClient(_ context.Context, opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	client, err := s.provider.CreateIdentityPoolClient(opts)
	if err != nil {
		return nil, errors.E(err, "creating an identity pool client", errors.Internal)
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
