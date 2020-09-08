package client

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
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
}
