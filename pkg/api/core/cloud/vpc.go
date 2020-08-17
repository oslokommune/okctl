// Package cloud implements the cloud layer
package cloud

import (
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

const (
	defaultTimeOut = 5
)

type vpc struct {
	provider v1alpha1.CloudProvider
}

// CreateCluster will use the cloud provider to create a cluster in the cloud
func (c *vpc) CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error) {
	b := cfn.New(components.NewVPCComposer(opts.RepoName, opts.Env, opts.Cidr, opts.Region))

	body, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build vpc cloud formation template", errors.Internal)
	}

	stackName := cfn.NewStackNamer().Vpc(opts.RepoName, opts.Env)

	r := cfn.NewRunner(stackName, body, c.provider)

	err = r.CreateIfNotExists(defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create vpc")
	}

	v, err := processOutputs(r, c.provider)
	if err != nil {
		return nil, errors.E(err, "failed to process vpc outputs")
	}

	v.StackName = stackName
	v.CloudFormationTemplate = body

	return v, nil
}

// DeleteVpc will use the cloud provider to delete a cluster in the cloud
func (c *vpc) DeleteVpc(opts api.DeleteVpcOpts) error {
	r := cfn.NewRunner(
		cfn.NewStackNamer().Vpc(opts.RepoName, opts.Env),
		nil,
		c.provider,
	)

	return r.Delete()
}

// NewVpcCloud returns a cloud provider for cluster
func NewVpcCloud(provider v1alpha1.CloudProvider) api.VpcCloud {
	return &vpc{
		provider: provider,
	}
}

// processOutputs extracts the outputs we are interested in from the cloud formation stack
func processOutputs(m *cfn.Runner, provider v1alpha1.CloudProvider) (*api.Vpc, error) {
	v := &api.Vpc{}

	err := m.Outputs(map[string]cfn.ProcessOutputFn{
		"PrivateSubnetIds": cfn.Subnets(provider, &v.PrivateSubnets),
		"PublicSubnetIds":  cfn.Subnets(provider, &v.PublicSubnets),
		"Vpc":              cfn.String(&v.ID),
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}
