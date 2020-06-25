// Package core implements the service layer
package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

type cluster struct {
	exe   api.ClusterExe
	cloud api.ClusterCloud
	store api.ClusterStore
}

// CreateCluster creates an EKS cluster and VPC
func (c *cluster) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	clusterConfig, err := c.cloud.CreateCluster(opts.AWSAccountID, opts.ClusterName, opts.Environment, opts.RepositoryName, opts.Cidr, opts.Region)
	if err != nil {
		return nil, err
	}

	err = c.exe.CreateCluster(clusterConfig)
	if err != nil {
		return nil, err
	}

	res := &api.Cluster{
		Environment:  opts.Environment,
		AWSAccountID: opts.AWSAccountID,
		Cidr:         opts.Cidr,
		Config:       clusterConfig,
	}

	err = c.store.SaveCluster(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteCluster deletes an EKS cluster and VPC
func (c *cluster) DeleteCluster(_ context.Context, opts api.ClusterDeleteOpts) error {
	clus, err := c.store.GetCluster(opts.Environment)
	if err != nil {
		return err
	}

	err = c.exe.DeleteCluster(clus.Config)
	if err != nil {
		return err
	}

	err = c.cloud.DeleteCluster(opts.Environment, opts.RepositoryName)
	if err != nil {
		return err
	}

	return c.store.DeleteCluster(opts.Environment)
}

// NewClusterService returns a service operator for the cluster operations
func NewClusterService(store api.ClusterStore, cloud api.ClusterCloud, exe api.ClusterExe) api.ClusterService {
	return &cluster{
		exe:   exe,
		cloud: cloud,
		store: store,
	}
}
