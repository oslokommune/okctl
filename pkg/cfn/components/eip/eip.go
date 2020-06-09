package eip

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type eip struct {
	name              string
	number            int
	gatewayAttachment cfn.Namer
}

func (i *eip) Resource() cloudformation.Resource {
	return &ec2.EIP{
		Domain: "vpc",
		AWSCloudFormationDependsOn: []string{
			i.gatewayAttachment.Name(),
		},
	}
}

func (i *eip) Name() string {
	return i.name
}

func (i *eip) Ref() string {
	return cloudformation.Ref(i.Name())
}

func New(number int, gatewayAttachment cfn.Namer) *eip {
	return &eip{
		name:              fmt.Sprintf("NatGatewayEIP%02d", number),
		number:            number,
		gatewayAttachment: gatewayAttachment,
	}
}
