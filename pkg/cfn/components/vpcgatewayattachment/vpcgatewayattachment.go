package vpcgatewayattachment

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type VPCGatewayAttachment struct {
	name string
	vpc  cfn.Referencer
	igw  cfn.Referencer
}

func New(vpc cfn.Referencer, igw cfn.Referencer) *VPCGatewayAttachment {
	return &VPCGatewayAttachment{
		name: "VPCGatewayAttachment",
		vpc:  vpc,
		igw:  igw,
	}
}

func (v *VPCGatewayAttachment) Resource() cloudformation.Resource {
	return &ec2.VPCGatewayAttachment{
		InternetGatewayId: v.igw.Ref(),
		VpcId:             v.vpc.Ref(),
	}
}

func (v *VPCGatewayAttachment) Name() string {
	return v.name
}

func (v *VPCGatewayAttachment) Ref() string {
	return cloudformation.Ref(v.Name())
}
