package securitygroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type securityGroup struct {
	name string
	vpc  cfn.Referencer
}

func (s *securityGroup) Resource() cloudformation.Resource {
	return &ec2.SecurityGroup{
		VpcId:            s.vpc.Ref(),
		GroupDescription: s.name,
	}
}

func (s *securityGroup) Name() string {
	return s.name
}

func (s *securityGroup) Ref() string {
	return cloudformation.Ref(s.Name())
}

func ControlPlane(vpc cfn.Referencer) *securityGroup {
	return &securityGroup{
		name: "ControlPlaneSecurityGroup",
		vpc:  vpc,
	}
}
