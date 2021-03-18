// Package vpcendpoint knows how to build a VPC endpoint
package vpcendpoint

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// VPCEndpoint contains the required state for
// building a cloud formation resource of a
// VPC endpoint
type VPCEndpoint struct {
	StoredName    string
	SecurityGroup cfn.Namer
	VpcID         string
	DBSubnetIDs   []string
	ServiceName   string
}

// NamedOutputs returns the resource outputs
func (e *VPCEndpoint) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(
		fmt.Sprintf("%sDnsEntries", e.Name()), cloudformation.GetAtt(e.Name(), "DnsEntries"),
	).NamedOutputs()
}

// Name returns the name of the cloud formation resource
func (e *VPCEndpoint) Name() string {
	return e.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (e *VPCEndpoint) Ref() string {
	return cloudformation.Ref(e.Name())
}

// Resource returns the cloud formation resource for a VPC endpoint
func (e *VPCEndpoint) Resource() cloudformation.Resource {
	return &ec2.VPCEndpoint{
		PrivateDnsEnabled: false,
		SecurityGroupIds: []string{
			cloudformation.GetAtt(e.SecurityGroup.Name(), "GroupId"),
		},
		ServiceName:     e.ServiceName,
		SubnetIds:       e.DBSubnetIDs,
		VpcEndpointType: "Interface",
		VpcId:           e.VpcID,
	}
}

// New returns an initialised VPC endpoint
// - https://docs.aws.amazon.com/vpc/latest/privatelink/vpc-endpoints.html
func New(
	resourceName string,
	securityGroup cfn.Namer,
	vpcID string,
	dbSubnets []string,
	serviceName string,
) *VPCEndpoint {
	return &VPCEndpoint{
		StoredName:    resourceName,
		SecurityGroup: securityGroup,
		VpcID:         vpcID,
		DBSubnetIDs:   dbSubnets,
		ServiceName:   serviceName,
	}
}

// NewSecretsManager returns an initialised vpc endpoint for the secrets manager
func NewSecretsManager(
	resourceName string,
	securityGroup cfn.Namer,
	vpcID string,
	dbSubnets []string,
) *VPCEndpoint {
	return New(
		resourceName,
		securityGroup,
		vpcID,
		dbSubnets,
		cloudformation.Sub(`com.amazonaws.${AWS::Region}.secretsmanager`),
	)
}
