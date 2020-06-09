package internetgateway

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
)

type InternetGateway struct {
	name string
}

func New() *InternetGateway {
	return &InternetGateway{
		name: "InternetGateway",
	}
}

func (i *InternetGateway) Resource() cloudformation.Resource {
	return &ec2.InternetGateway{}
}

func (i *InternetGateway) Name() string {
	return i.name
}

func (i *InternetGateway) Ref() string {
	return cloudformation.Ref(i.Name())
}
