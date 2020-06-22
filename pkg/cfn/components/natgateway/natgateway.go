package natgateway

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type NatGateway struct {
	StoredName        string
	Number            int
	EIP               cfn.Namer
	Subnet            cfn.NameReferencer
	GatewayAttachment cfn.Namer
}

func (n *NatGateway) Resource() cloudformation.Resource {
	return &ec2.NatGateway{
		AllocationId: cloudformation.GetAtt(n.EIP.Name(), "AllocationId"),
		SubnetId:     n.Subnet.Ref(),
		AWSCloudFormationDependsOn: []string{
			n.EIP.Name(),
			n.Subnet.Name(),
			n.GatewayAttachment.Name(),
		},
	}
}

func (n *NatGateway) Name() string {
	return n.StoredName
}

func (n *NatGateway) Ref() string {
	return cloudformation.Ref(n.Name())
}

// New returns a cloud formation VPC NAT gateway attachment
//
// Attaches an internet gateway, or a virtual private gateway to a VPC,
// enabling connectivity between the internet and the VPC.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpc-gateway-attachment.html
func New(number int, gatewayAttachment cfn.Namer, eip cfn.Namer, subnet cfn.NameReferencer) *NatGateway {
	return &NatGateway{
		StoredName:        fmt.Sprintf("NatGateway%02d", number),
		Number:            0,
		EIP:               eip,
		Subnet:            subnet,
		GatewayAttachment: gatewayAttachment,
	}
}
