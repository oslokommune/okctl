// Package aliasrecordset provides cloud formation for a
// Route53 recordset with an Alias to an AWS resource, such
// as a cloud front distribution
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-recordset.html
package aliasrecordset

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/route53"
)

// AliasRecordSet contains all state for creating a record set
// cloud formation template
type AliasRecordSet struct {
	StoredName        string
	AliasDNS          string
	AliasHostedZoneID string
	Domain            string
	HostedZoneID      string
}

// Resource returns the cloudformation resource
func (s *AliasRecordSet) Resource() cloudformation.Resource {
	return &route53.RecordSet{
		AliasTarget: &route53.RecordSet_AliasTarget{
			DNSName:      s.AliasDNS,
			HostedZoneId: s.AliasHostedZoneID,
		},
		HostedZoneId: s.HostedZoneID,
		Name:         s.Domain,
		Type:         "A",
	}
}

// Name is the logical id of the resource
func (s *AliasRecordSet) Name() string {
	return s.StoredName
}

// New returns an initialised record set creator
func New(name, aliasDNS, aliasHostedZoneID, domain, hostedZoneID string) *AliasRecordSet {
	return &AliasRecordSet{
		StoredName:        fmt.Sprintf("AliasRecordSet%s", name),
		AliasDNS:          aliasDNS,
		AliasHostedZoneID: aliasHostedZoneID,
		Domain:            domain,
		HostedZoneID:      hostedZoneID,
	}
}
