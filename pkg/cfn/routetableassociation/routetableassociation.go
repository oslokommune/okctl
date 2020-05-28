package routetableassociation

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type association struct {
	name       string
	subnet     cfn.Referencer
	routeTable cfn.Referencer
}

func (a *association) Resource() cloudformation.Resource {
	return &ec2.SubnetRouteTableAssociation{
		RouteTableId: a.routeTable.Ref(),
		SubnetId:     a.subnet.Ref(),
	}
}

func (a *association) Name() string {
	return a.name
}

func (a *association) Ref() string {
	return cloudformation.Ref(a.Name())
}

func newAssociation(number int, t string, subnet cfn.Referencer, routeTable cfn.Referencer) *association {
	return &association{
		name:       fmt.Sprintf("%sSubnet%02dRouteTableAssociation", t, number),
		subnet:     subnet,
		routeTable: routeTable,
	}
}

func NewPublic(number int, subnet cfn.Referencer, routeTable cfn.Referencer) *association {
	return newAssociation(number, "Public", subnet, routeTable)
}

func NewPrivate(number int, subnet cfn.Referencer, routeTable cfn.Referencer) *association {
	return newAssociation(number, "Private", subnet, routeTable)
}
