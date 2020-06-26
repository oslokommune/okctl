// Package cloud implements the cloud layer
package cloud

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/oslokommune/okctl/pkg/cfn/process"
)

const (
	defaultTimeOut = 5
)

type cluster struct {
	provider v1alpha1.CloudProvider
}

// CreateCluster will use the cloud provider to create a cluster in the cloud
func (c *cluster) CreateCluster(awsAccountID, clusterName, env, repoName, cidr, region string) (*v1alpha1.ClusterConfig, error) {
	m := manager.New(c.provider).
		WithBuilder(vpc.New(repoName, env, cidr, region))

	err := m.CreateIfNotExists(defaultTimeOut)
	if err != nil {
		return nil, err
	}

	clusterConfig, err := clusterConfig(clusterName, region, cidr, awsAccountID, m, c.provider)
	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}

// DeleteCluster will use the cloud provider to delete a cluster in the cloud
func (c *cluster) DeleteCluster(env, repoName string) error {
	return manager.New(c.provider).Delete(vpc.StackName(repoName, env))
}

// NewCluster returns a cloud provider for cluster
func NewCluster(provider v1alpha1.CloudProvider) api.ClusterCloud {
	return &cluster{
		provider: provider,
	}
}

// clusterConfig knows how to populate an eksctl cluster config after vpc is created
func clusterConfig(name, region, cidr, awsAccountID string, m *manager.Manager, provider v1alpha1.CloudProvider) (*v1alpha1.ClusterConfig, error) {
	clusterConfig := v1alpha1.NewClusterConfig()

	clusterConfig.Metadata.Name = name
	clusterConfig.Metadata.Region = region
	clusterConfig.VPC.CIDR = cidr
	clusterConfig.IAM.FargatePodExecutionRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(awsAccountID)
	clusterConfig.IAM.ServiceRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(awsAccountID)

	return clusterConfig, m.Outputs(map[string]manager.ProcessOutputFn{
		"PrivateSubnetIds": process.Subnets(provider, clusterConfig.VPC.Subnets.Private),
		"PublicSubnetIds":  process.Subnets(provider, clusterConfig.VPC.Subnets.Public),
		"Vpc":              process.String(&clusterConfig.VPC.ID),
	})
}
