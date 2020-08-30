// Package certificate creates components for an ACM public certificate
package certificate

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/certificatemanager"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// Certificate contains state for building a cloud formation resource
type Certificate struct {
	StoredName   string
	FQDN         string
	HostedZoneID string
}

// NamedOutputs returns the named outputs
func (c Certificate) NamedOutputs() map[string]map[string]interface{} {
	return cfn.NewValue(c.Name(), c.Ref()).NamedOutputs()
}

// Resource returns the cloud formation resource
func (c *Certificate) Resource() cloudformation.Resource {
	return &certificatemanager.Certificate{
		DomainName: c.FQDN,
		DomainValidationOptions: []certificatemanager.Certificate_DomainValidationOption{
			{
				DomainName:   c.FQDN,
				HostedZoneId: c.HostedZoneID,
			},
		},
		ValidationMethod: "DNS",
	}
}

// Name returns the logical identifier
func (c *Certificate) Name() string {
	return c.StoredName
}

// Ref returns an aws intrinsic ref to this resource
func (c *Certificate) Ref() string {
	return cloudformation.Ref(c.Name())
}

// New initialises a new certificate
func New(fqdn, hostedZoneID string) *Certificate {
	return &Certificate{
		StoredName:   "PublicCertificate",
		FQDN:         fqdn,
		HostedZoneID: hostedZoneID,
	}
}
