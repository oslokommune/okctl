package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Certificate contains the state for a certificate
type Certificate struct {
	Repository             string
	Environment            string
	FQDN                   string
	Domain                 string
	HostedZoneID           string
	CertificateARN         string
	StackName              string
	CloudFormationTemplate []byte
}

// CreateCertificateOpts contains the input required for creating a certificate
type CreateCertificateOpts struct {
	Repository   string
	Environment  string
	FQDN         string
	Domain       string
	HostedZoneID string
}

// Validate the input
func (o CreateCertificateOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Repository, validation.Required),
		validation.Field(&o.Environment, validation.Required),
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
}

// CertificateStore defines the storage operations
type CertificateStore interface {
	SaveCertificate(certificate *Certificate) error
}
