// Package hostedzone knows how to create cloud formation for a hosted zone
package hostedzone

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/route53"
	"github.com/oslokommune/okctl/pkg/cfn"
	"gopkg.in/yaml.v3"
)

// HostedZone contains the state for creating a cloud formation resources and outputs
type HostedZone struct {
	StoredName string
	FQDN       string
	Comment    string
}

// NamedOutputs returns the named outputs
func (h *HostedZone) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValueMap().
		Add(cfn.NewValue(h.Name(), h.Ref())).
		// This doesn't work, so we need to patch this after the fact, this
		// is because cloudformation.Join doesn't support Fn::Join with an
		// Fn::GetAtt.
		Add(cfn.NewValue("NameServers", "I need to be patched.")).
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

// Patcher struct contains what we want
// to patch
type Patcher struct {
	AWSTemplateFormatVersion string `yaml:"AWSTemplateFormatVersion"`
	Outputs                  struct {
		NameServers struct {
			Value interface{} `yaml:"Value"`
		} `yaml:"NameServers"`
		PublicHostedZone struct {
			Value struct {
				Ref string `yaml:"Ref"`
			} `yaml:"Value"`
		} `yaml:"PublicHostedZone"`
	} `yaml:"Outputs"`
	Resources struct {
		PublicHostedZone struct {
			Properties struct {
				HostedZoneConfig struct {
					Comment string `yaml:"Comment"`
				} `yaml:"HostedZoneConfig"`
				Name string `yaml:"Name"`
			} `yaml:"Properties"`
			Type string `yaml:"Type"`
		} `yaml:"PublicHostedZone"`
	} `yaml:"Resources"`
}

// PatchYAML the template body, so it is valid
// cloud formation
func PatchYAML(templateBody []byte) ([]byte, error) {
	getAtt := map[string][]string{
		"Fn::GetAtt": {
			"PublicHostedZone",
			"NameServers",
		},
	}
	list := []interface{}{
		",",
		getAtt,
	}
	join := map[string]interface{}{
		"Fn::Join": list,
	}

	patched := &Patcher{}

	err := yaml.Unmarshal(templateBody, patched)
	if err != nil {
		return nil, err
	}

	patched.Outputs.NameServers.Value = join

	d, err := yaml.Marshal(patched)
	if err != nil {
		return nil, err
	}

	return d, nil
}
