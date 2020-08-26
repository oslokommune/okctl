// Package hostedzone knows how to create cloud formation for a hosted zone
package hostedzone

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/route53"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// HostedZone contains the state for creating a cloud formation resources and outputs
type HostedZone struct {
	StoredName string
	FQDN       string
	Comment    string
}

// NamedOutputs returns the named outputs
func (h *HostedZone) NamedOutputs() map[string]map[string]interface{} {
	return cfn.NewValueMap().
		Add(cfn.NewValue(h.Name(), h.Ref())).
		Add(cfn.NewValue("NameServers", cloudformation.Join(",", []string{h.GetAtt("NameServers")}))).
		NamedOutputs()
}

// Resource returns the cloud formation resource for a HostedZone
func (h *HostedZone) Resource() cloudformation.Resource {
	return &route53.HostedZone{
		HostedZoneConfig: &route53.HostedZone_HostedZoneConfig{
			Comment: h.Comment,
		},
		Name: h.FQDN,
	}
}

// Name returns the name of the resource
func (h *HostedZone) Name() string {
	return h.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (h *HostedZone) Ref() string {
	return cloudformation.Ref(h.Name())
}

// GetAtt returns a cloud formation intrinsic GetAtt to an attribute
func (h *HostedZone) GetAtt(attribute string) string {
	return cloudformation.GetAtt(h.Name(), attribute)
}

// New returns an initialised hosted zone
func New(fqdn, comment string) *HostedZone {
	return &HostedZone{
		StoredName: "PublicHostedZone",
		FQDN:       fqdn,
		Comment:    comment,
	}
}
