// Package routetable knows how to create cloud formation for a route table
package routetable

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// RouteTable stores required state for creating a
// cloud formation route table resource
type RouteTable struct {
	StoredName string
	Number     int
	VPC        cfn.Referencer
}

// Resource returns the cloud formation route table
func (p *RouteTable) Resource() cloudformation.Resource {
	return &ec2.RouteTable{
		VpcId: p.VPC.Ref(),
	}
}

// Name returns the name of the resource
func (p *RouteTable) Name() string {
	return p.StoredName
}

// Ref return a cloud formation intrinsic ref to the resource
func (p *RouteTable) Ref() string {
	return cloudformation.Ref(p.Name())
}

// NewPrivate returns a route table for the RouteTable
// subnets
//
// Specifies a route table for a specified VPC.
// After you create a route table, you can add routes
// and associate the table with a subnet.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route-table.html
func NewPrivate(number int, vpc cfn.Referencer) *RouteTable {
	return &RouteTable{
		StoredName: fmt.Sprintf("PrivateRouteTable%02d", number),
		Number:     number,
		VPC:        vpc,
	}
}

// NewPublic returns a route table for the Public
// subnets
//
// Specifies a route table for a specified VPC.
// After you create a route table, you can add routes
// and associate the table with a subnet.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route-table.html
func NewPublic(vpc cfn.Referencer) *RouteTable {
	return &RouteTable{
		StoredName: "PublicRouteTable",
		VPC:        vpc,
	}
}
