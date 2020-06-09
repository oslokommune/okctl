package routetable

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type private struct {
	name   string
	number int
	vpc    cfn.Referencer
}

func NewPrivate(number int, vpc cfn.Referencer) *private {
	return &private{
		name:   fmt.Sprintf("PrivateRouteTable%02d", number),
		number: number,
		vpc:    vpc,
	}
}

func (p *private) Resource() cloudformation.Resource {
	return &ec2.RouteTable{
		VpcId: p.vpc.Ref(),
	}
}

func (p *private) Name() string {
	return p.name
}

func (p *private) Ref() string {
	return cloudformation.Ref(p.Name())
}

type public struct {
	name string
	vpc  cfn.Referencer
}

func (p *public) Resource() cloudformation.Resource {
	return &ec2.RouteTable{
		VpcId: p.vpc.Ref(),
	}
}

func (p *public) Name() string {
	return p.name
}

func (p *public) Ref() string {
	return cloudformation.Ref(p.Name())
}

func NewPublic(vpc cfn.Referencer) *public {
	return &public{
		name: "PublicRouteTable",
		vpc:  vpc,
	}
}
