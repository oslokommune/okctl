// Package aws implements the cloud layer
package aws

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

const (
	defaultTimeOut = 5
)

type vpcCloudProvider struct {
	provider v1alpha1.CloudProvider
}

// CreateVpc will use the cloud provider to create a cluster in the cloud
func (c *vpcCloudProvider) CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error) {
	var b *cfn.Builder
	if opts.Minimal {
		b = cfn.New(components.NewMinimalVPCComposer(opts.ID.ClusterName, opts.Cidr, opts.ID.Region))
	} else {
		b = cfn.New(components.NewVPCComposer(opts.ID.ClusterName, opts.Cidr, opts.ID.Region))
	}

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf(constant.BuildCloudFormationTemplateError, err)
	}

	stackName := cfn.NewStackNamer().Vpc(opts.ID.ClusterName)

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(opts.ID.ClusterName, stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, fmt.Errorf(constant.CreateVPCStackError, err)
	}

	v := &api.Vpc{
		ID:                     opts.ID,
		StackName:              stackName,
		CloudFormationTemplate: template,
		Cidr:                   opts.Cidr,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"Vpc":                     cfn.String(&v.VpcID),
		"PrivateSubnetIds":        cfn.Subnets(c.provider, &v.PrivateSubnets),
		"PublicSubnetIds":         cfn.Subnets(c.provider, &v.PublicSubnets),
		"DatabaseSubnetIds":       cfn.Subnets(c.provider, &v.DatabaseSubnets),
		"DatabaseSubnetGroupName": cfn.String(&v.DatabaseSubnetsGroupName),
	})
	if err != nil {
		return nil, fmt.Errorf(constant.ProcessStackOutputError, err)
	}

	return v, nil
}

// DeleteVpc will use the cloud provider to delete a cluster in the cloud
func (c *vpcCloudProvider) DeleteVpc(opts api.DeleteVpcOpts) error {
	return cfn.NewRunner(c.provider).Delete(cfn.NewStackNamer().Vpc(opts.ID.ClusterName))
}

// NewVpcCloud returns a cloud provider for cluster
func NewVpcCloud(provider v1alpha1.CloudProvider) api.VpcCloudProvider {
	return &vpcCloudProvider{
		provider: provider,
	}
}
