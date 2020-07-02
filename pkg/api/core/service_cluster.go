// Package core implements the service layer
package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type cluster struct {
	exe   api.ClusterExe
	cloud api.ClusterCloud
	store api.ClusterStore
}

const (
	msgFailedToCreateCluster = "failed to create cluster"
	msgFailedToDeleteCluster = "failed to delete cluster"
)

// CreateCluster creates an EKS cluster and VPC
func (c *cluster) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate create cluster input", errors.Invalid)
	}

	clusterConfig, err := c.cloud.CreateCluster(opts.AWSAccountID, opts.ClusterName, opts.Environment, opts.RepositoryName, opts.Cidr, opts.Region)
	if err != nil {
		return nil, errors.E(err, msgFailedToCreateCluster)
	}

	err = c.exe.CreateCluster(clusterConfig)
	if err != nil {
		return nil, errors.E(err, msgFailedToCreateCluster)
	}

	res := &api.Cluster{
		Environment:  opts.Environment,
		AWSAccountID: opts.AWSAccountID,
		Cidr:         opts.Cidr,
		Config:       clusterConfig,
	}

	err = c.store.SaveCluster(res)
	if err != nil {
		return nil, errors.E(err, msgFailedToCreateCluster)
	}

	return res, nil
}

// DeleteCluster deletes an EKS cluster and VPC
func (c *cluster) DeleteCluster(_ context.Context, opts api.ClusterDeleteOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "failed to validate delete cluster inputs", errors.Invalid)
	}

	clus, err := c.store.GetCluster(opts.Environment)
	if err != nil {
		return errors.E(err, "failed to get cluster from config", errors.Invalid)
	}

	err = c.exe.DeleteCluster(clus.Config)
	if err != nil {
		return errors.E(err, msgFailedToDeleteCluster)
	}

	err = c.cloud.DeleteCluster(opts.Environment, opts.RepositoryName)
	if err != nil {
		return errors.E(err, msgFailedToDeleteCluster)
	}

	return errors.E(c.store.DeleteCluster(opts.Environment), msgFailedToDeleteCluster)
}

// NewClusterService returns a service operator for the cluster operations
func NewClusterService(store api.ClusterStore, cloud api.ClusterCloud, exe api.ClusterExe) api.ClusterService {
	return &cluster{
		exe:   exe,
		cloud: cloud,
		store: store,
	}
}
