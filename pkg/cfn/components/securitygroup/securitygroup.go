// Package securitygroup knows how to create cloud formation for security groups
package securitygroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/builder/output"
)

// SecurityGroup stores state required for creating a
// cloud formation security group
type SecurityGroup struct {
	StoredName string
	VPC        cfn.Referencer
}

// NamedOutputs returns the outputs commonly used by other stacks or components
func (s *SecurityGroup) NamedOutputs() map[string]map[string]interface{} {
	return output.NewValue(s.Name(), s.Ref()).NamedOutputs()
}

// Resource returns the cloud formation resource for creating a SG
func (s *SecurityGroup) Resource() cloudformation.Resource {
	return &ec2.SecurityGroup{
		VpcId:            s.VPC.Ref(),
		GroupDescription: s.StoredName,
	}
}

// Name returns the name of the cloud formation resource
func (s *SecurityGroup) Name() string {
	return s.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (s *SecurityGroup) Ref() string {
	return cloudformation.Ref(s.Name())
}

// ControlPlane creates an EKS control plane security group
func ControlPlane(vpc cfn.Referencer) *SecurityGroup {
	return &SecurityGroup{
		StoredName: "ControlPlaneSecurityGroup",
		VPC:        vpc,
	}
}
