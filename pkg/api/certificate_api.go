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

// Validate the input
func (o CreateCertificateOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.FQDN, validation.Required),
		validation.Field(&o.Domain, validation.Required),
		validation.Field(&o.HostedZoneID, validation.Required),
	)
}

// DeleteCertificateOpts contains input required to delete a certificate
type DeleteCertificateOpts struct {
	ID     ID
	Domain string
}

// Validate the deletion request inputs
func (o DeleteCertificateOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Domain, validation.Required),
	)
}

// DeleteCognitoCertificateOpts contains the inputs
type DeleteCognitoCertificateOpts struct {
	ID     ID
	Domain string
}

// Validate the deletion request inputs
func (o DeleteCognitoCertificateOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Domain, validation.Required),
	)
}

// CertificateService defines the service layer operations
type CertificateService interface {
	CreateCertificate(ctx context.Context, opts CreateCertificateOpts) (*Certificate, error)
	DeleteCertificate(ctx context.Context, opts DeleteCertificateOpts) error
	DeleteCognitoCertificate(ctx context.Context, opts DeleteCognitoCertificateOpts) error
}

// CertificateCloudProvider defines the cloud interaction
type CertificateCloudProvider interface {
	CreateCertificate(opts CreateCertificateOpts) (*Certificate, error)
	DeleteCertificate(opts DeleteCertificateOpts) error
	DeleteCognitoCertificate(opts DeleteCognitoCertificateOpts) error
}
