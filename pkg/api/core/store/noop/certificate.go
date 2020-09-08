package noop

import "github.com/oslokommune/okctl/pkg/api"

type certificateStore struct{}

func (s *certificateStore) SaveCertificate(_ *api.Certificate) error {
	return nil
}

// NewCertificateStore returns a no operation store
func NewCertificateStore() api.CertificateStore {
	return &certificateStore{}
}
