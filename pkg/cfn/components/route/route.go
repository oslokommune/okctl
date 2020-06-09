package route

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type private struct {
	name              string
	number            int
	routeTable        cfn.Referencer
	natGateway        cfn.NameReferencer
	gatewayAttachment cfn.Namer
}

func (p *private) Ref() string {
	return cloudformation.Ref(p.Name())
}

func (p *private) Resource() cloudformation.Resource {
	return &ec2.Route{
		DestinationCidrBlock: "0.0.0.0/0",
		NatGatewayId:         p.natGateway.Ref(),
		RouteTableId:         p.routeTable.Ref(),
		AWSCloudFormationDependsOn: []string{
			p.natGateway.Name(),
			p.gatewayAttachment.Name(),
		},
	}
}

func (p *private) Name() string {
	return p.name
}

func NewPrivate(number int, gatewayAttachment cfn.Namer, routeTable cfn.Referencer, natGateway cfn.NameReferencer) *private {
	return &private{
		name:              fmt.Sprintf("PrivateRoute%02d", number),
		number:            number,
		routeTable:        routeTable,
		natGateway:        natGateway,
		gatewayAttachment: gatewayAttachment,
	}
}

type public struct {
	name              string
	routeTable        cfn.Referencer
	gatewayAttachment cfn.Namer
	internetGateway   cfn.Referencer
}

func (p *public) Resource() cloudformation.Resource {
	return &ec2.Route{
		DestinationCidrBlock: "0.0.0.0/0",
		GatewayId:            p.internetGateway.Ref(),
		RouteTableId:         p.routeTable.Ref(),
		AWSCloudFormationDependsOn: []string{
			p.gatewayAttachment.Name(),
		},
	}
}

func (p *public) Name() string {
	return p.name
}

func (p *public) Ref() string {
	return cloudformation.Ref(p.Name())
}

func NewPublic(gatewayAttachment cfn.Namer, routeTable cfn.Referencer, internetGateway cfn.Referencer) *public {
	return &public{
		name:              "PublicRoute",
		routeTable:        routeTable,
		gatewayAttachment: gatewayAttachment,
		internetGateway:   internetGateway,
	}
}
