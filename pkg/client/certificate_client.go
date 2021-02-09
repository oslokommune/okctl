package client

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

// We are shadowing some interfaces for now, but
// this is probably not sustainable.

// CertificateService orchestrates the creation of a certificate
type CertificateService interface {
	api.CertificateService
}

// CertificateAPI defines the api interactions
type CertificateAPI interface {
	api.CertificateCloudProvider
}

// CertificateStore defines the storage operations
type CertificateStore interface {
	SaveCertificate(certificate *api.Certificate) (*store.Report, error)
	GetCertificate(domain string) (*api.Certificate, error)
	RemoveCertificate(domain string) (*store.Report, error)
}

// CertificateState defines the state layer
type CertificateState interface {
	SaveCertificate(certificate *api.Certificate) (*store.Report, error)
	GetCertificate(domain string) state.Certificate
	RemoveCertificate(domain string) (*store.Report, error)
}

// CertificateReport defines the report layer
type CertificateReport interface {
	SaveCertificate(certificate *api.Certificate, reports []*store.Report) error
	RemoveCertificate(domain string, reports []*store.Report) error
}
