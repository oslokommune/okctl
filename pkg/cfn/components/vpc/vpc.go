package vpc

import (
	"fmt"
	"net"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/awslabs/goformation/v4/cloudformation/tags"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type VPC struct {
	name    string
	cluster cfn.Namer
	block   *net.IPNet
}

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

func (v *VPC) Name() string {
	return v.name
}

func (v *VPC) Ref() string {
	return cloudformation.Ref(v.Name())
}

func New(cluster cfn.Namer, cidr *net.IPNet) *VPC {
	return &VPC{
		name:    "VPC",
		cluster: cluster,
		block:   cidr,
	}
}
