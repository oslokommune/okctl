// Package cloud implements the cloud layer
package cloud

import (
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	vpcBuilder "github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/oslokommune/okctl/pkg/cfn/process"
)

const (
	defaultTimeOut = 5
)

type vpc struct {
	provider v1alpha1.CloudProvider
}

// CreateCluster will use the cloud provider to create a cluster in the cloud
func (c *vpc) CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error) {
	builder := vpcBuilder.New(opts.RepoName, opts.Env, opts.Cidr, opts.Region)

	body, err := builder.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build vpc cloud formation template", errors.Internal)
	}

	m := manager.New(builder.StackName(), body, c.provider)

	err = m.CreateIfNotExists(defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create vpc")
	}

	v, err := processOutputs(m, c.provider)
	if err != nil {
		return nil, errors.E(err, "failed to process vpc outputs")
	}

	v.StackName = vpcBuilder.StackName(opts.RepoName, opts.Env)
	v.CloudFormationTemplate = body

	return v, nil
}

// DeleteVpc will use the cloud provider to delete a cluster in the cloud
func (c *vpc) DeleteVpc(opts api.DeleteVpcOpts) error {
	return manager.New(vpcBuilder.StackName(opts.RepoName, opts.Env), nil, c.provider).Delete()
}

// NewVpcCloud returns a cloud provider for cluster
func NewVpcCloud(provider v1alpha1.CloudProvider) api.VpcCloud {
	return &vpc{
		provider: provider,
	}
}

// processOutputs extracts the outputs we are interested in from the cloud formation stack
func processOutputs(m *manager.Manager, provider v1alpha1.CloudProvider) (*api.Vpc, error) {
	v := &api.Vpc{}

	err := m.Outputs(map[string]manager.ProcessOutputFn{
		"PrivateSubnetIds": process.Subnets(provider, &v.PrivateSubnets),
		"PublicSubnetIds":  process.Subnets(provider, &v.PublicSubnets),
		"Vpc":              process.String(&v.ID),
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}
