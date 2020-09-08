// Package core implements the service layer
package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type clusterService struct {
	store api.ClusterStore
	run   api.ClusterRun
}

func (s *clusterService) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate inputs", errors.Invalid)
	}

	cluster, err := s.run.CreateCluster(opts)
	if err != nil {
		return nil, errors.E(err, "failed to create cluster", errors.Internal)
	}

	err = s.store.SaveCluster(cluster)
	if err != nil {
		return nil, errors.E(err, "failed to save cluster", errors.IO)
	}

	return cluster, nil
}

func (s *clusterService) DeleteCluster(_ context.Context, opts api.ClusterDeleteOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "failed to validate inputs", errors.Invalid)
	}

	err = s.run.DeleteCluster(opts)
	if err != nil {
		return errors.E(err, "failed to delete cluster", errors.Internal)
	}

	err = s.store.DeleteCluster(opts.ID)
	if err != nil {
		return errors.E(err, "failed to remove cluster", errors.IO)
	}

	return nil
}

// NewClusterService returns a service operator for the clusterService operations
func NewClusterService(store api.ClusterStore, run api.ClusterRun) api.ClusterService {
	return &clusterService{
		store: store,
		run:   run,
	}
}
