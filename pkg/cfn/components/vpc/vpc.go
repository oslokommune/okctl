// Package vpc knows how to create a cloud formation VPC
package vpc

import (
	"fmt"
	"net"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/awslabs/goformation/v4/cloudformation/tags"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// VPC stores the state for creating a cloud formation VPC
type VPC struct {
	name    string
	cluster cfn.Namer
	block   *net.IPNet
}

// NamedOutputs returns the commonly used named outputs of a VPC
func (v *VPC) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(v.Name(), v.Ref()).NamedOutputs()
}

// Resource returns the cloud formation resource of the VPC
func (v *VPC) Resource() cloudformation.Resource {
	t := []tags.Tag{
		{
			Key:   fmt.Sprintf("kubernetes.io/cluster/%s", v.cluster.Name()),
			Value: "shared",
		},
	}

	return &ec2.VPC{
		CidrBlock:          v.block.String(),
		EnableDnsHostnames: true,
		EnableDnsSupport:   true,
		Tags:               t,
	}
}

// Name returns the name of the resource
func (v *VPC) Name() string {
	return v.name
}

// Ref returns a cloud formation intrinsic ref to the resource
func (v *VPC) Ref() string {
	return cloudformation.Ref(v.Name())
}

// New returns a new VPC cloud formation resource
func New(cluster cfn.Namer, cidr *net.IPNet) *VPC {
	return &VPC{
		name:    "Vpc",
		cluster: cluster,
		block:   cidr,
	}
}
