// Package internetgateway knows how to create cloud formation for a IGW
package internetgateway

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
)

// InternetGateway stores state required for creating
// a cloud formation IGW
type InternetGateway struct {
	name string
}

// New returns a new IGW
func New() *InternetGateway {
	return &InternetGateway{
		name: "InternetGateway",
	}
}

// Resource returns the cloud formation resource for an IGW
func (i *InternetGateway) Resource() cloudformation.Resource {
	return &ec2.InternetGateway{}
}

// Name returns the name of the resource
func (i *InternetGateway) Name() string {
	return i.name
}

// Ref returns a cloud formation intrinsic ref to the resource
func (i *InternetGateway) Ref() string {
	return cloudformation.Ref(i.Name())
}
