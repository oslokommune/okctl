package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
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

	cfg, err := clusterconfig.New(&clusterconfig.Args{
		ClusterName:            opts.ClusterName,
		PermissionsBoundaryARN: v1alpha1.PermissionsBoundaryARN(opts.AwsAccountID),
		PrivateSubnets:         vpc.PrivateSubnets,
		PublicSubnets:          vpc.PublicSubnets,
		Region:                 opts.Region,
		VpcCidr:                opts.Cidr,
		VpcID:                  vpc.ID,
	})
	if err != nil {
		return nil, errors.E(err, "failed to create cluster config")
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
