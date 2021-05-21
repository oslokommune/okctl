package storm

import (
	"errors"
	"time"

	"github.com/oslokommune/okctl/pkg/breeze"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type certificateState struct {
	node breeze.Client
}

// Certificate contains storm compatible state
type Certificate struct {
	Metadata `storm:"inline"`

	ID                     ID
	FQDN                   string
	Domain                 string `storm:"unique"`
	HostedZoneID           string
	ARN                    string
	StackName              string
	CloudFormationTemplate string
}

// NewCertificate constructs a storm compatible  Certificate
func NewCertificate(c *client.Certificate, meta Metadata) *Certificate {
	return &Certificate{
		Metadata:               meta,
		ID:                     NewID(c.ID),
		FQDN:                   c.FQDN,
		Domain:                 c.Domain,
		HostedZoneID:           c.HostedZoneID,
		ARN:                    c.ARN,
		StackName:              c.StackName,
		CloudFormationTemplate: string(c.CloudFormationTemplate),
	}
}

// Convert a Certificate to *client.Certificate
func (c *Certificate) Convert() *client.Certificate {
	return &client.Certificate{
		ID:                     c.ID.Convert(),
		FQDN:                   c.FQDN,
		Domain:                 c.Domain,
		HostedZoneID:           c.HostedZoneID,
		ARN:                    c.ARN,
		StackName:              c.StackName,
		CloudFormationTemplate: []byte(c.CloudFormationTemplate),
	}
}

func (c *certificateState) SaveCertificate(certificate *client.Certificate) error {
	existing, err := c.getCertificate(certificate.Domain)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return c.node.Save(NewCertificate(certificate, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return c.node.Save(NewCertificate(certificate, existing.Metadata))
}

func (c *certificateState) getCertificate(domain string) (*Certificate, error) {
	cert := &Certificate{}

	err := c.node.One("Domain", domain, cert)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (c *certificateState) GetCertificate(domain string) (*client.Certificate, error) {
	cert, err := c.getCertificate(domain)
	if err != nil {
		return nil, err
	}

	return cert.Convert(), nil
}

func (c *certificateState) RemoveCertificate(domain string) error {
	cert := &Certificate{}

	err := c.node.One("Domain", domain, cert)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return c.node.DeleteStruct(cert)
}

// NewCertificateState returns an initialised state store
func NewCertificateState(node breeze.Client) client.CertificateState {
	return &certificateState{
		node: node,
	}
}
