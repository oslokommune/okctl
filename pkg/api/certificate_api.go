package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Certificate contains the state for a certificate
type Certificate struct {
	ID                     ID
	FQDN                   string
	Domain                 string
	HostedZoneID           string
	CertificateARN         string
	StackName              string
	CloudFormationTemplate []byte
}

// CreateCertificateOpts contains the input required for creating a certificate
type CreateCertificateOpts struct {
	ID           ID
	FQDN         string
	Domain       string
	HostedZoneID string
}

// DeleteCertificateOpts contains input required to delete a certificate
type DeleteCertificateOpts struct {
	Domain string
}

// Validate the input
func (o CreateCertificateOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.FQDN, validation.Required),
		validation.Field(&o.Domain, validation.Required),
		validation.Field(&o.HostedZoneID, validation.Required),
	)
}

// CertificateService defines the service layer operations
type CertificateService interface {
	CreateCertificate(ctx context.Context, opts CreateCertificateOpts) (*Certificate, error)
}

// CertificateCloudProvider defines the cloud interaction
type CertificateCloudProvider interface {
	CreateCertificate(opts CreateCertificateOpts) (*Certificate, error)
	DeleteCertificate(opts DeleteCertificateOpts) error
}

// CertificateStore defines the storage operations
type CertificateStore interface {
	SaveCertificate(certificate *Certificate) error
}
