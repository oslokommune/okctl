// Package routetableassociation knows how to create cloud formation for a route table association
package routetableassociation

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// RouteTableAssociation stores the state of the cloud formation resource
type RouteTableAssociation struct {
	StoredName string
	Subnet     cfn.Referencer
	RouteTable cfn.Referencer
}

// Resource returns a cloud formation resource of the route table association
func (a *RouteTableAssociation) Resource() cloudformation.Resource {
	return &ec2.SubnetRouteTableAssociation{
		RouteTableId: a.RouteTable.Ref(),
		SubnetId:     a.Subnet.Ref(),
	}
}

// Name returns the name of the resource
func (a *RouteTableAssociation) Name() string {
	return a.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (a *RouteTableAssociation) Ref() string {
	return cloudformation.Ref(a.Name())
}

func newAssociation(number int, t string, subnet cfn.Referencer, routeTable cfn.Referencer) *RouteTableAssociation {
	return &RouteTableAssociation{
		StoredName: fmt.Sprintf("%sSubnet%02dRouteTableAssociation", t, number),
		Subnet:     subnet,
		RouteTable: routeTable,
	}
}

// NewPublic returns a public subnet route table association
//
// Associates a subnet with a route table. This association causes
// traffic originating from the subnet to be routed according to
// the routes in the route table
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-subnet-route-table-assoc.html
func NewPublic(number int, subnet cfn.Referencer, routeTable cfn.Referencer) *RouteTableAssociation {
	return newAssociation(number, "Public", subnet, routeTable)
}

// NewPrivate returns a private subnet route table association
//
// Associates a subnet with a route table. This association causes
// traffic originating from the subnet to be routed according to
// the routes in the route table
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-subnet-route-table-assoc.html
func NewPrivate(number int, subnet cfn.Referencer, routeTable cfn.Referencer) *RouteTableAssociation {
	return newAssociation(number, "Private", subnet, routeTable)
}
