package natgateway

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type natGateway struct {
	name              string
	number            int
	eip               cfn.Namer
	subnet            cfn.NameReferencer
	gatewayAttachment cfn.Namer
}

func (n *natGateway) Resource() cloudformation.Resource {
	return &ec2.NatGateway{
		AllocationId: cloudformation.GetAtt(n.eip.Name(), "AllocationId"),
		SubnetId:     n.subnet.Ref(),
		AWSCloudFormationDependsOn: []string{
			n.eip.Name(),
			n.subnet.Name(),
			n.gatewayAttachment.Name(),
		},
	}
}

func (n *natGateway) Name() string {
	return n.name
}

func (n *natGateway) Ref() string {
	return cloudformation.Ref(n.Name())
}

func New(number int, gatewayAttachment cfn.Namer, eip cfn.Namer, subnet cfn.NameReferencer) *natGateway {
	return &natGateway{
		name:              fmt.Sprintf("NatGateway%02d", number),
		number:            0,
		eip:               eip,
		subnet:            subnet,
		gatewayAttachment: gatewayAttachment,
	}
}
