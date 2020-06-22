// Package route knows how to create cloud formation for an EC2 VPC route
package route

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// Route stores required state for creating a cloud formation Route
type Route struct {
	Type              string
	Number            int
	StoredName        string
	RouteTable        cfn.Referencer
	GatewayAttachment cfn.Namer
	InternetGateway   cfn.Referencer
	NatGateway        cfn.NameReferencer
}

// Resource returns the cloud formation route resource
func (r *Route) Resource() cloudformation.Resource {
	route := &ec2.Route{
		DestinationCidrBlock: "0.0.0.0/0",
		RouteTableId:         r.RouteTable.Ref(),
		AWSCloudFormationDependsOn: []string{
			r.GatewayAttachment.Name(),
		},
	}

	switch r.Type {
	case "public":
		route.GatewayId = r.InternetGateway.Ref()
	case "private":
		route.NatGatewayId = r.NatGateway.Ref()
		route.AWSCloudFormationDependsOn = append(route.AWSCloudFormationDependsOn, r.NatGateway.Name())
	}

	return route
}

// Name returns the name of the resource
func (r *Route) Name() string {
	return r.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (r *Route) Ref() string {
	return cloudformation.Ref(r.Name())
}

// NewPrivate creates a private EC2 route, routing traffic to the NATGW
//
// Specifies a route in a route table within a VPC.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route.html
func NewPrivate(number int, gatewayAttachment cfn.Namer, routeTable cfn.Referencer, natGateway cfn.NameReferencer) *Route {
	return &Route{
		Type:              "private",
		StoredName:        fmt.Sprintf("PrivateRoute%02d", number),
		Number:            number,
		RouteTable:        routeTable,
		NatGateway:        natGateway,
		GatewayAttachment: gatewayAttachment,
	}
}

// NewPublic creates a Route EC2 route, routing traffic to the IGW
//
// Specifies a route in a route table within a VPC.
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-route.html
func NewPublic(gatewayAttachment cfn.Namer, routeTable cfn.Referencer, internetGateway cfn.Referencer) *Route {
	return &Route{
		Type:              "public",
		StoredName:        "PublicRoute",
		RouteTable:        routeTable,
		GatewayAttachment: gatewayAttachment,
		InternetGateway:   internetGateway,
	}
}
