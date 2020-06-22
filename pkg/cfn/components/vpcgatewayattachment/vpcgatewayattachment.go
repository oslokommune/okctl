// Package vpcgatewayattachment knows how to create cloud formation for a vpc gateway attachment
package vpcgatewayattachment

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// VPCGatewayAttachment stores the state required for creating a
// cloud formation gateway attachment
type VPCGatewayAttachment struct {
	name string
	vpc  cfn.Referencer
	igw  cfn.Referencer
}

// New creates a new VPC gateway attachment
//
// Attaches an internet gateway, or a virtual private
// gateway to a VPC, enabling connectivity between the
// internet and the VPC.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpc-gateway-attachment.html
func New(vpc cfn.Referencer, igw cfn.Referencer) *VPCGatewayAttachment {
	return &VPCGatewayAttachment{
		name: "VPCGatewayAttachment",
		vpc:  vpc,
		igw:  igw,
	}
}

// Resource returns the cloud formation resource
func (v *VPCGatewayAttachment) Resource() cloudformation.Resource {
	return &ec2.VPCGatewayAttachment{
		InternetGatewayId: v.igw.Ref(),
		VpcId:             v.vpc.Ref(),
	}
}

// Name returns the name of the resource
func (v *VPCGatewayAttachment) Name() string {
	return v.name
}

// Ref returns a cloud formation intrinsic ref to the resource
func (v *VPCGatewayAttachment) Ref() string {
	return cloudformation.Ref(v.Name())
}
