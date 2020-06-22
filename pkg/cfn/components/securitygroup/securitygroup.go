package securitygroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/builder/output"
)

type SecurityGroup struct {
	StoredName string
	VPC        cfn.Referencer
}

func (s *SecurityGroup) NamedOutputs() map[string]map[string]interface{} {
	return output.NewValue(s.Name(), s.Ref()).NamedOutputs()
}

func (s *SecurityGroup) Resource() cloudformation.Resource {
	return &ec2.SecurityGroup{
		VpcId:            s.VPC.Ref(),
		GroupDescription: s.StoredName,
	}
}

func (s *SecurityGroup) Name() string {
	return s.StoredName
}

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
