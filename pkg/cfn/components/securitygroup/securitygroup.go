// Package securitygroup knows how to create cloud formation for security groups
package securitygroup

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// SecurityGroup stores state required for creating a
// cloud formation security group
type SecurityGroup struct {
	StoredName string
	Group      *ec2.SecurityGroup
}

// NamedOutputs returns the outputs commonly used by other stacks or components
func (s *SecurityGroup) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValueMap().
		Add(cfn.NewValue(s.Name(), s.Ref())).
		Add(cfn.NewValue(fmt.Sprintf("%s-GroupId", s.Name()), cloudformation.GetAtt(s.Name(), "GroupId"))).
		NamedOutputs()
}

// Resource returns the cloud formation resource for creating a SG
func (s *SecurityGroup) Resource() cloudformation.Resource {
	return s.Group
}

// Name returns the name of the cloud formation resource
func (s *SecurityGroup) Name() string {
	return s.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (s *SecurityGroup) Ref() string {
	return cloudformation.Ref(s.Name())
}

const (
	postgresPort = 5432
)

// NewPostgresOutgoing returns an initialised security group
// that allows outgoing traffic from the pod or node to the
// postgres subnets on the postgres port
func NewPostgresOutgoing(resourceName, vpcID string, cidrs []string) *SecurityGroup {
	egresses := make([]ec2.SecurityGroup_Egress, len(cidrs))

	for i, cidr := range cidrs {
		egresses[i] = ec2.SecurityGroup_Egress{
			CidrIp:     cidr,
			FromPort:   postgresPort,
			IpProtocol: "tcp",
			ToPort:     postgresPort,
		}
	}

	return &SecurityGroup{
		StoredName: resourceName,
		Group: &ec2.SecurityGroup{
			GroupDescription:    "RDS Postgres Outgoing Security Group",
			GroupName:           resourceName,
			SecurityGroupEgress: egresses,
			VpcId:               vpcID,
		},
	}
}

// NewPostgresIncoming returns an initialised security group that
// allows incoming traffic to the postgres database instance
func NewPostgresIncoming(resourceName, vpcID string, source cfn.Namer) *SecurityGroup {
	return &SecurityGroup{
		StoredName: resourceName,
		Group: &ec2.SecurityGroup{
			GroupDescription:    "RDS Postgres Incoming Security Group",
			GroupName:           resourceName,
			SecurityGroupEgress: nil,
			SecurityGroupIngress: []ec2.SecurityGroup_Ingress{
				{
					FromPort:              postgresPort,
					IpProtocol:            "tcp",
					SourceSecurityGroupId: cloudformation.GetAtt(source.Name(), "GroupId"),
					ToPort:                postgresPort,
				},
			},
			VpcId: vpcID,
		},
	}
}
