package routetable

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type Private struct {
	StoredName string
	Number     int
	VPC        cfn.Referencer
}

// NewPrivate returns a route table for the Private
// subnets
//
// Specifies a route table for a specified VPC.
// After you create a route table, you can add routes
// and associate the table with a subnet.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route-table.html
func NewPrivate(number int, vpc cfn.Referencer) *Private {
	return &Private{
		StoredName: fmt.Sprintf("PrivateRouteTable%02d", number),
		Number:     number,
		VPC:        vpc,
	}
}

func (p *Private) Resource() cloudformation.Resource {
	return &ec2.RouteTable{
		VpcId: p.VPC.Ref(),
	}
}

func (p *Private) Name() string {
	return p.StoredName
}

func (p *Private) Ref() string {
	return cloudformation.Ref(p.Name())
}

type Public struct {
	StoredName string
	VPC        cfn.Referencer
}

func (p *Public) Resource() cloudformation.Resource {
	return &ec2.RouteTable{
		VpcId: p.VPC.Ref(),
	}
}

func (p *Public) Name() string {
	return p.StoredName
}

func (p *Public) Ref() string {
	return cloudformation.Ref(p.Name())
}

// NewPublic returns a route table for the Public
// subnets
//
// Specifies a route table for a specified VPC.
// After you create a route table, you can add routes
// and associate the table with a subnet.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route-table.html
func NewPublic(vpc cfn.Referencer) *Public {
	return &Public{
		StoredName: "PublicRouteTable",
		VPC:        vpc,
	}
}
