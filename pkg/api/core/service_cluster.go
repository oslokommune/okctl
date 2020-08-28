// Package core implements the service layer
package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type cluster struct {
	run             api.ClusterRun
	store           api.ClusterStore
	kubeConfigStore api.KubeConfigStore
}

// CreateCluster creates an EKS cluster and VPC
func (c *cluster) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate create cluster input", errors.Invalid)
	}

	kubeConfigPath, err := c.kubeConfigStore.CreateKubeConfig()
	if err != nil {
		return nil, errors.E(err, "failed to create kubeconfig", errors.IO)
	}

	cfg, err := clusterconfig.New(&clusterconfig.Args{
		ClusterName:            opts.ClusterName,
		PermissionsBoundaryARN: v1alpha1.PermissionsBoundaryARN(opts.AWSAccountID),
		PrivateSubnets:         opts.VpcPrivateSubnets,
		PublicSubnets:          opts.VpcPublicSubnets,
		Region:                 opts.Region,
		VpcCidr:                opts.Cidr,
		VpcID:                  opts.VpcID,
	})
	if err != nil {
		return nil, errors.E(err, "failed to create cluster config", errors.Internal)
	}

	err = c.run.CreateCluster(kubeConfigPath, cfg)
	if err != nil {
		return nil, errors.E(err, "failed to create cluster", errors.Internal)
	}

	res := &api.Cluster{
		Environment:    opts.Environment,
		AWSAccountID:   opts.AWSAccountID,
		Cidr:           opts.Cidr,
		ClusterName:    opts.ClusterName,
		RepositoryName: opts.RepositoryName,
		Region:         opts.Region,
		Config:         cfg,
	}

	err = c.store.SaveCluster(res)
	if err != nil {
		return nil, errors.E(err, "failed to save cluster", errors.IO)
	}

	return res, nil
}

// DeleteClusterConfig deletes an EKS cluster
func (c *cluster) DeleteCluster(_ context.Context, opts api.ClusterDeleteOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "failed to validate delete cluster inputs", errors.Invalid)
	}

	err = c.run.DeleteCluster(opts.ClusterName)
	if err != nil {
		return errors.E(err, "failed to delete cluster", errors.Internal)
	}

	err = c.store.DeleteCluster(opts.Environment)
	if err != nil {
		return errors.E(err, "failed to delete cluster", errors.Internal)
	}

	return nil
}

// NewClusterService returns a service operator for the cluster operations
func NewClusterService(store api.ClusterStore, kubeConfigStore api.KubeConfigStore, run api.ClusterRun) api.ClusterService {
	return &cluster{
		run:             run,
		store:           store,
		kubeConfigStore: kubeConfigStore,
	}
}
