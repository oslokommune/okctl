package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

type clusterConfig struct {
	store    api.ClusterConfigStore
	vpcStore api.VpcStore
}

// CreateClusterConfig implements the business logic for creating a cluster config
func (c *clusterConfig) CreateClusterConfig(_ context.Context, opts api.CreateClusterConfigOpts) (*api.ClusterConfig, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate create cluster config input")
	}

	vpc, err := c.vpcStore.GetVpc()
	if err != nil {
		return nil, errors.E(err, "failed to retrieve stored vpc state")
	}

	cfg := api.NewClusterConfig()

	cfg.Metadata.Name = opts.ClusterName
	cfg.Metadata.Region = opts.Region
	cfg.VPC.ID = vpc.ID
	cfg.VPC.CIDR = opts.Cidr
	cfg.IAM.FargatePodExecutionRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(opts.AwsAccountID)
	cfg.IAM.ServiceRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(opts.AwsAccountID)

	for _, p := range vpc.PublicSubnets {
		cfg.VPC.Subnets.Public[p.AvailabilityZone] = api.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	for _, p := range vpc.PrivateSubnets {
		cfg.VPC.Subnets.Private[p.AvailabilityZone] = api.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	err = c.store.SaveClusterConfig(cfg)
	if err != nil {
		return nil, errors.E(err, "failed to store cluster config")
	}

	return cfg, nil
}

// NewClusterConfigService returns an instantiated cluster config service
func NewClusterConfigService(store api.ClusterConfigStore, vpcStore api.VpcStore) api.ClusterConfigService {
	return &clusterConfig{
		store:    store,
		vpcStore: vpcStore,
	}
}
