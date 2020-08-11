// Package natgateway knows how to create cloud formation for a NATGW
package natgateway

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// NatGateway stores required state for creating
// a cloud formation NATGW
type NatGateway struct {
	StoredName        string
	Number            int
	EIP               cfn.Namer
	PublicSubnet      cfn.NameReferencer
	GatewayAttachment cfn.Namer
}

// Resource returns a cloud formation resource for creating a NATGW
func (n *NatGateway) Resource() cloudformation.Resource {
	return &ec2.NatGateway{
		AllocationId: cloudformation.GetAtt(n.EIP.Name(), "AllocationId"),
		SubnetId:     n.PublicSubnet.Ref(),
		AWSCloudFormationDependsOn: []string{
			n.EIP.Name(),
			n.PublicSubnet.Name(),
			n.GatewayAttachment.Name(),
		},
	}
}

// Name returns the name of the resource
func (n *NatGateway) Name() string {
	return n.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (n *NatGateway) Ref() string {
	return cloudformation.Ref(n.Name())
}

// New returns a cloud formation VPC NAT gateway attachment
//
// Attaches an internet gateway, or a virtual private gateway to a VPC,
// enabling connectivity between the internet and the VPC.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpc-gateway-attachment.html
func New(number int, gatewayAttachment cfn.Namer, eip cfn.Namer, publicSubnet cfn.NameReferencer) *NatGateway {
	return &NatGateway{
		StoredName:        fmt.Sprintf("NatGateway%02d", number),
		Number:            0,
		EIP:               eip,
		PublicSubnet:      publicSubnet,
		GatewayAttachment: gatewayAttachment,
	}
}
