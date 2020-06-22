// Package eip provides a simplified interface for creating cloud
// formation resources
package eip

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// EIP stores the state for creating a cloud formation EIP resource
type EIP struct {
	StoredName        string
	Number            int
	GatewayAttachment cfn.Namer
}

// Resource returns the cloud formation resource for an EIP
func (i *EIP) Resource() cloudformation.Resource {
	return &ec2.EIP{
		Domain: "vpc",
		AWSCloudFormationDependsOn: []string{
			i.GatewayAttachment.Name(),
		},
	}
}

// Name returns the name of the resource
func (i *EIP) Name() string {
	return i.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
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
