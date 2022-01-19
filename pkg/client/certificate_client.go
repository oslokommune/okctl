package client

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/api"
)

// Certificate represents an AWS ACM certificate
// after it has been created
type Certificate struct {
	ID                     api.ID
	FQDN                   string
	Domain                 string
	HostedZoneID           string
	ARN                    string
	StackName              string
	CloudFormationTemplate []byte
}

// CreateCertificateOpts contains the input required for creating a certificate
type CreateCertificateOpts struct {
	ID           api.ID
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
	ID     api.ID
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
	ID     api.ID
	Domain string
}

// Validate the deletion request inputs
func (o DeleteCognitoCertificateOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Domain, validation.Required),
	)
}

// CertificateService orchestrates the creation of a certificate
type CertificateService interface {
	CreateCertificate(ctx context.Context, opts CreateCertificateOpts) (*Certificate, error)
	DeleteCertificate(ctx context.Context, opts DeleteCertificateOpts) error
	DeleteCognitoCertificate(ctx context.Context, opts DeleteCognitoCertificateOpts) error
}

// CertificateAPI defines the api interactions
type CertificateAPI interface {
	CreateCertificate(opts api.CreateCertificateOpts) (*api.Certificate, error)
	DeleteCertificate(opts api.DeleteCertificateOpts) error
	DeleteCognitoCertificate(opts api.DeleteCognitoCertificateOpts) error
}

// CertificateState defines the state layer
type CertificateState interface {
	SaveCertificate(certificate *Certificate) error
	GetCertificate(domain string) (*Certificate, error)
	HasCertificate(domain string) (bool, error)
	RemoveCertificate(domain string) error
}
