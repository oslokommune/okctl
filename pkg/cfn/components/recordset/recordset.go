// Package recordset provides cloud formation for a Route53 recordset
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-recordset.html
package recordset

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/route53"
)

// RecordSet contains all state for creating a record set
// cloud formation template
type RecordSet struct {
	StoredName   string
	Domain       string
	IP           string
	HostedZoneID string
}

// Resource returns the cloudformation resource
func (s *RecordSet) Resource() cloudformation.Resource {
	return &route53.RecordSet{
		HostedZoneId: s.HostedZoneID,
		ResourceRecords: []string{
			s.IP,
		},
		Name: s.Domain,
		Type: "A",
	}
}

// Name is the logical id of the resource
func (s *RecordSet) Name() string {
	return s.StoredName
}

// New returns an initialised record set creator
func New(name, ip, domain, hostedZoneID string) *RecordSet {
	return &RecordSet{
		StoredName:   fmt.Sprintf("RecordSet%s", name),
		Domain:       domain,
		IP:           ip,
		HostedZoneID: hostedZoneID,
	}
}
