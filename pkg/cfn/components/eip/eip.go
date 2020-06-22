package eip

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type EIP struct {
	StoredName        string
	Number            int
	GatewayAttachment cfn.Namer
}

func (i *EIP) Resource() cloudformation.Resource {
	return &ec2.EIP{
		Domain: "vpc",
		AWSCloudFormationDependsOn: []string{
			i.GatewayAttachment.Name(),
		},
	}
}

func (i *EIP) Name() string {
	return i.StoredName
}

func (i *EIP) Ref() string {
	return cloudformation.Ref(i.Name())
}

// New creates an EIP for use with the NAT GW
//
// Specifies an Elastic IP (EIP) address
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-eip.html
func New(number int, gatewayAttachment cfn.Namer) *EIP {
	return &EIP{
		StoredName:        fmt.Sprintf("NatGatewayEIP%02d", number),
		Number:            number,
		GatewayAttachment: gatewayAttachment,
	}
}
