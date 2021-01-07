// Package aws implements the cloud layer
package aws

import (
	"github.com/mishudark/errors"
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

// CreateCluster will use the cloud provider to create a cluster in the cloud
func (c *vpcCloudProvider) CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error) {
	var b *cfn.Builder
	if opts.Minimal {
		b = cfn.New(components.NewMinimalVPCComposer(opts.ID.Repository, opts.ID.Environment, opts.Cidr, opts.ID.Region))
	} else {
		b = cfn.New(components.NewVPCComposer(opts.ID.Repository, opts.ID.Environment, opts.Cidr, opts.ID.Region))
	}

	template, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template", errors.Internal)
	}

	stackName := cfn.NewStackNamer().Vpc(opts.ID.Repository, opts.ID.Environment)

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create vpc")
	}

	v := &api.Vpc{
		ID:                     opts.ID,
		StackName:              stackName,
		CloudFormationTemplate: template,
		Cidr:                   opts.Cidr,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"PrivateSubnetIds": cfn.Subnets(c.provider, &v.PrivateSubnets),
		"PublicSubnetIds":  cfn.Subnets(c.provider, &v.PublicSubnets),
		"Vpc":              cfn.String(&v.VpcID),
	})
	if err != nil {
		return nil, errors.E(err, "failed to process outputs")
	}

	return v, nil
}

// DeleteVpc will use the cloud provider to delete a cluster in the cloud
func (c *vpcCloudProvider) DeleteVpc(opts api.DeleteVpcOpts) error {
	return cfn.NewRunner(c.provider).Delete(cfn.NewStackNamer().Vpc(opts.ID.Repository, opts.ID.Environment))
}

// NewVpcCloud returns a cloud provider for cluster
func NewVpcCloud(provider v1alpha1.CloudProvider) api.VpcCloudProvider {
	return &vpcCloudProvider{
		provider: provider,
	}
}
