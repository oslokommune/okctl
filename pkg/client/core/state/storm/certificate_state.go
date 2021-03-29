package storm

import (
	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type certificateState struct {
	node stormpkg.Node
}

// Certificate contains storm compatible state
type Certificate struct {
	Metadata `storm:"inline"`

	ID                     ID
	FQDN                   string
	Domain                 string `storm:"unique,index"`
	HostedZoneID           string
	ARN                    string
	StackName              string
	CloudFormationTemplate []byte
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
		CloudFormationTemplate: c.CloudFormationTemplate,
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
		CloudFormationTemplate: c.CloudFormationTemplate,
	}
}

func (c *certificateState) SaveCertificate(certificate *client.Certificate) error {
	return c.node.Save(NewCertificate(certificate, NewMetadata()))
}

func (c *certificateState) GetCertificate(domain string) (*client.Certificate, error) {
	var cert Certificate

	err := c.node.One("Domain", domain, &cert)
	if err != nil {
		return nil, err
	}

	return cert.Convert(), nil
}

func (c *certificateState) RemoveCertificate(domain string) error {
	var cert Certificate

	err := c.node.One("Domain", domain, &cert)
	if err != nil {
		return err
	}

	return c.node.DeleteStruct(&cert)
}

// NewCertificateState returns an initialised state store
func NewCertificateState(node stormpkg.Node) client.CertificateState {
	return &certificateState{
		node: node,
	}
}
